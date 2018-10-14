package rhodeisland

import (
	"errors"
	. "github.com/int8/gopoker"
)

// RIGameState - RhodeIslandGameState
type RIGameState struct {
	round         Round
	parent        *RIGameState
	causingAction Action
	table         *Table
	actors        map[ActorId]Actor
	nextToMove    ActorId
	terminal      bool
}

func (state *RIGameState) Act(action Action) GameState {
	switch state.actors[state.nextToMove].(type) {
	case *Chance:
		return state.actAsChance(action)
	case *Player:
		return state.actAsPlayer(action)
	}
	return nil
}

func (state *RIGameState) Actions() []Action {

	switch state.actors[state.nextToMove].(type) {
	case *Chance:
		return state.chanceActions(state.chanceActor())
	case *Player:
		return state.playerActions(state.playerActor(state.nextToMove))
	}
	return nil
}

func (state *RIGameState) IsChance() bool {
	return state.nextToMove == ChanceId
}

func (state *RIGameState) IsTerminal() bool {
	return state.terminal
}

func (state *RIGameState) Parent() GameState {
	return state.parent
}

func (state *RIGameState) CurrentActor() Actor {
	return state.actors[state.nextToMove]
}

//TODO: test it carefully
func (state *RIGameState) Evaluate() float32 {
	currentActor := state.playerActor(state.CurrentActor().GetId())
	currentActorOpponent := state.playerActor(-state.CurrentActor().GetId())
	if state.IsTerminal() {
		if state.causingAction.Name() == Fold {
			currentActor.UpdateStack(currentActor.stack + state.table.Pot)
			return float32(state.nextToMove) * state.table.Pot / 2
		}
		currentActorHandValueVector := currentActor.EvaluateHand(state.table)
		currentActorOpponentHandValueVector := currentActorOpponent.EvaluateHand(state.table)
		for i := range currentActorHandValueVector {
			if currentActorHandValueVector[i] == currentActorOpponentHandValueVector[i] {
				continue
			}
			if currentActorHandValueVector[i] > currentActorOpponentHandValueVector[i] {
				currentActor.UpdateStack(currentActor.stack + state.table.Pot)
				return float32(currentActor.GetId()) * state.table.Pot / 2
			} else {
				currentActorOpponent.UpdateStack(currentActorOpponent.stack + state.table.Pot)
				return float32(currentActorOpponent.GetId()) * state.table.Pot / 2
			}
		}
		state.playerActor(-state.CurrentActor().GetId()).UpdateStack(state.table.Pot / 2)
		state.playerActor(state.CurrentActor().GetId()).UpdateStack(state.table.Pot / 2)
		return 0.0
	}
	panic(errors.New("RIGameState is not terminal"))
}

func (state *RIGameState) InformationSet() InformationSet {

	privateCardName := byte(state.playerActor(state.nextToMove).card.Name)
	privateCardSuit := byte(state.playerActor(state.nextToMove).card.Suit)
	flopCardName := byte(NoCardName)
	flopCardSuit := byte(NoCardSuit)
	turnCardName := byte(NoCardName)
	turnCardSuit := byte(NoCardSuit)

	if len(state.table.Cards) > 0 {
		flopCardName = byte(state.table.Cards[0].Name)
		flopCardSuit = byte(state.table.Cards[0].Suit)
	}

	if len(state.table.Cards) > 1 {
		turnCardName = byte(state.table.Cards[1].Name)
		turnCardSuit = byte(state.table.Cards[1].Suit)
	}
	informationSet := [InformationSetSize]byte{privateCardName, privateCardSuit, flopCardName, flopCardSuit, turnCardName, turnCardSuit}
	// there is no more than 50 actions overall
	currentState := state
	for i := 6; currentState.round != Start; i++ {
		informationSet[i] = byte(currentState.causingAction.Name())
		currentState = currentState.parent
		if currentState == nil {
			break
		}
	}
	return InformationSet(informationSet)
}

func (state *RIGameState) stack(actor ActorId) float32 {
	return state.actors[actor].(*Player).stack
}

func (state *RIGameState) actAsChance(action Action) GameState {
	var child *RIGameState
	if action.Name() == DealPublicCards {
		child = state.dealPublicCard(action.(DealPublicCardAction).Card)
	}

	if action.Name() == DealPrivateCards {
		child = state.dealPrivateCards(action.(DealPrivateCardsAction).CardA, action.(DealPrivateCardsAction).CardB)
	}
	return child
}

func (state *RIGameState) actAsPlayer(action Action) GameState {

	var child *RIGameState

	if !actionInSlice(action, state.Actions()) {
		panic("action not available")
	}
	actor := state.playerActor(state.nextToMove)
	betSize := state.betSize()

	defer func() {
		if action.Name() == Call || action.Name() == Bet {
			child.playerActor(actor.GetId()).PlaceBet(child.table, betSize)
		}
		if action.Name() == Raise {
			child.playerActor(actor.GetId()).PlaceBet(child.table, 2*betSize)
		}
		if action.Name() == Fold {
			child.playerActor(-actor.GetId()).PlaceBet(child.table, -betSize)
		}
	}()

	if action.Name() == Fold || (state.round == Turn && (action.Name() == Call || (action.Name() == Check && state.causingAction.Name() == Check))) {
		child = createChild(state, state.round, action, actor.Opponent(), true)
		return child
	}

	if action.Name() == Call || (action.Name() == Check && state.causingAction.Name() == Check) {
		child = createChild(state, state.round, action, ChanceId, false)
		return child
	}

	child = createChild(state, state.round, action, actor.Opponent(), false)
	return child

}

func (state *RIGameState) betSize() float32 {
	if state.round < Flop {
		return PreFlopBetSize
	}
	return PostFlopBetSize
}

func root(playerA *Player, playerB *Player) *RIGameState {
	chance := &Chance{id: ChanceId, deck: CreateFullDeck(true)}

	actors := map[ActorId]Actor{PlayerA: playerA, PlayerB: playerB, ChanceId: chance}
	table := &Table{Pot: 0, Cards: []Card{}}

	return &RIGameState{round: Start, table: table,
		actors: actors, nextToMove: ChanceId, causingAction: nil}
}

func createChild(blueprint *RIGameState, round Round, Action Action, nextToMove ActorId, terminal bool) *RIGameState {
	child := RIGameState{round: round,
		parent: blueprint, causingAction: Action, terminal: terminal,
		table: blueprint.table.Clone(), actors: cloneActorsMap(blueprint.actors), nextToMove: nextToMove}
	return &child
}

func (state *RIGameState) dealPublicCard(card *Card) *RIGameState {

	child := createChild(state, state.round.NextRound(), DealPublicCardAction{card}, PlayerA, false)
	// important to deal using child deck / not current chance deck
	child.table.DropPublicCard(card)
	child.actors[ChanceId].(*Chance).deck.RemoveCard(card)
	return child
}

func (state *RIGameState) dealPrivateCards(cardA *Card, cardB *Card) *RIGameState {

	child := createChild(state, state.round.NextRound(), DealPrivateCardsAction{cardA, cardB}, PlayerA, false)
	// important to deal using child deck / not current chance deck
	child.playerActor(PlayerA).PlaceBet(child.table, Ante)
	child.actors[PlayerB].(*Player).PlaceBet(child.table, Ante)
	child.playerActor(PlayerA).CollectPrivateCard(cardA)
	child.actors[PlayerB].(*Player).CollectPrivateCard(cardB)
	child.actors[ChanceId].(*Chance).deck.RemoveCard(cardA)
	child.actors[ChanceId].(*Chance).deck.RemoveCard(cardB)

	return child
}

func (state *RIGameState) chanceActions(chance *Chance) []Action {
	if state.round == Start {
		deckSize := int(chance.deck.CardsLeft())
		actions := make([]Action, deckSize*(deckSize-1))
		i := 0
		remainingCards := chance.deck.RemainingCards()
		for _, cardA := range remainingCards {
			for _, cardB := range remainingCards {
				{
					if cardA != cardB {
						actions[i] = DealPrivateCardsAction{cardA, cardB}
						i++
					}
				}
			}
		}
		return actions
	}

	actions := make([]Action, chance.deck.CardsLeft())
	remainingCards := chance.deck.RemainingCards()
	for i, card := range remainingCards {
		actions[i] = DealPublicCardAction{card}
	}
	return actions
}

func (state *RIGameState) playerActions(player *Player) []Action {

	if state.causingAction.Name() == Fold {
		player.actions = []Action{}
		return player.actions
	}

	betSize := state.betSize()

	opponentStack := state.stack(player.Opponent())

	allowedToBet := (player.stack >= betSize) && (opponentStack >= betSize)
	allowedToRaise := (player.stack >= 2*betSize) && (opponentStack >= 2*betSize)

	// whenever betting roung is over (CALL OR CHECK->CHECK)
	bettingRoundEnded := state.causingAction.Name() == Call || (state.causingAction.Name() == Check && state.parent.causingAction.Name() == Check)
	if bettingRoundEnded {
		player.actions = []Action{}
		return player.actions
	}

	// single check implies BET or CHECK
	if state.causingAction.Name() == Check && state.parent.causingAction.Name() != Check {
		player.actions = []Action{CheckAction}
		if allowedToBet {
			player.actions = append(player.actions, BetAction)
		}
		return player.actions
	}

	// RAISE/BET, you can CALL FOLD or RAISE (unless there has been 6 prior raises - 3 for each player)
	if state.causingAction.Name() == Bet || state.causingAction.Name() == Raise {
		player.actions = []Action{CallAction, FoldAction}
		priorRaisesInCurrentRound := countPriorRaisesPerRound(state, state.round)
		if priorRaisesInCurrentRound < MaxRaises && allowedToRaise {
			player.actions = append(player.actions, RaiseAction)
		}
		return player.actions
	}

	if state.causingAction.Name() == DealPrivateCards || state.causingAction.Name() == DealPublicCards {
		player.actions = []Action{CheckAction}
		if allowedToBet {
			player.actions = append(player.actions, BetAction)
		}
		return player.actions
	}
	panic(errors.New("Code not reachable."))
}

func (state *RIGameState) playerActor(id ActorId) *Player {
	return state.actors[id].(*Player)
}

func (state *RIGameState) chanceActor() *Chance {
	return state.actors[ChanceId].(*Chance)
}
