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
	actors        map[ActorID]Actor
	nextToMove    ActorID
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

func (state *RIGameState) Evaluate() float32 {
	actor := state.playerActor(state.CurrentActor().GetID())
	opponent := state.playerActor(-state.CurrentActor().GetID())
	if state.IsTerminal() {
		if state.causingAction.Name() == Fold {
			actor.UpdateStack(actor.Stack + state.table.Pot)
			return float32(actor.GetID()) * state.table.Pot / 2
		}
		actorHandVector := actor.EvaluateHand(state.table)
		opponentHandVector := opponent.EvaluateHand(state.table)
		for i := range actorHandVector {
			if actorHandVector[i] == opponentHandVector[i] {
				continue
			}
			if actorHandVector[i] > opponentHandVector[i] {
				actor.UpdateStack(actor.Stack + state.table.Pot)
				return float32(actor.GetID()) * state.table.Pot / 2
			}
			opponent.UpdateStack(opponent.Stack + state.table.Pot)
			return float32(opponent.GetID()) * state.table.Pot / 2

		}
		state.playerActor(opponent.GetID()).UpdateStack(state.table.Pot / 2)
		state.playerActor(state.CurrentActor().GetID()).UpdateStack(state.table.Pot / 2)
		return 0.0
	}
	panic(errors.New("RIGameState is not terminal"))
}

func (state *RIGameState) InformationSet() InformationSet {

	privateCard := Card{state.playerActor(state.nextToMove).Card.Symbol, state.playerActor(state.nextToMove).Card.Suit}
	flopCard, turnCard := NoCard, NoCard

	if len(state.table.Cards) > 0 {
		flopCard = Card{state.table.Cards[0].Symbol, state.table.Cards[0].Suit}
	}

	if len(state.table.Cards) > 1 {
		turnCard = Card{state.table.Cards[1].Symbol, state.table.Cards[1].Suit}
	}

	infSet := [InformationSetSize]bool{
		privateCard.Symbol[0], privateCard.Symbol[1], privateCard.Symbol[2], privateCard.Symbol[3],
		privateCard.Suit[0], privateCard.Suit[1], privateCard.Suit[2],
		flopCard.Symbol[0], flopCard.Symbol[1], flopCard.Symbol[2], flopCard.Symbol[3],
		flopCard.Suit[0], flopCard.Suit[1], flopCard.Suit[2],
		turnCard.Symbol[0], turnCard.Symbol[1], turnCard.Symbol[2], turnCard.Symbol[3],
		turnCard.Suit[0], turnCard.Suit[1], turnCard.Suit[2],
	}

	currentState := state
	for i := 21; currentState.round != Start; i += 3 {
		actionName := currentState.causingAction.Name()
		infSet[i] = actionName[0]
		infSet[i+1] = actionName[1]
		infSet[i+2] = actionName[2]

		currentState = currentState.parent
		if currentState == nil {
			break
		}
	}
	return InformationSet(infSet)
}

func (state *RIGameState) stack(id ActorID) float32 {
	return state.actors[id].(*Player).Stack
}

func (state *RIGameState) actAsChance(action Action) GameState {
	var c *RIGameState
	if action.Name() == DealPublicCards {
		c = state.dealPublicCard(action.(DealPublicCardAction).Card)
	}

	if action.Name() == DealPrivateCards {
		c = state.dealPrivateCards(action.(DealPrivateCardsAction).CardA, action.(DealPrivateCardsAction).CardB)
	}
	return c
}

func (state *RIGameState) actAsPlayer(action Action) GameState {

	var c *RIGameState

	if !actionInSlice(action, state.Actions()) {
		panic("action not available")
	}
	actor := state.playerActor(state.nextToMove)
	betSize := state.betSize()

	defer func() {
		if action.Name() == Call || action.Name() == Bet {
			c.playerActor(actor.GetID()).PlaceBet(c.table, betSize)
		}
		if action.Name() == Raise {
			c.playerActor(actor.GetID()).PlaceBet(c.table, 2*betSize)
		}
		if action.Name() == Fold {
			c.playerActor(-actor.GetID()).PlaceBet(c.table, -betSize)
		}
	}()

	if action.Name() == Fold || (state.round == Turn && (action.Name() == Call || (action.Name() == Check && state.causingAction.Name() == Check))) {
		c = createChild(state, state.round, action, actor.Opponent(), true)
		return c
	}

	if action.Name() == Call || (action.Name() == Check && state.causingAction.Name() == Check) {
		c = createChild(state, state.round, action, ChanceId, false)
		return c
	}

	c = createChild(state, state.round, action, actor.Opponent(), false)
	return c

}

func (state *RIGameState) betSize() float32 {
	if state.round < Flop {
		return PreFlopBetSize
	}
	return PostFlopBetSize
}

func Root(playerA *Player, playerB *Player, deck Deck) *RIGameState {
	chance := &Chance{id: ChanceId, deck: deck}

	actors := map[ActorID]Actor{PlayerA: playerA, PlayerB: playerB, ChanceId: chance}
	table := &Table{Pot: 0, Cards: []Card{}}
	return &RIGameState{round: Start, table: table,
		actors: actors, nextToMove: ChanceId, causingAction: nil}
}

func createChild(blueprint *RIGameState, round Round, action Action, nextToMove ActorID, terminal bool) *RIGameState {
	c := RIGameState{round: round,
		parent: blueprint, causingAction: action, terminal: terminal,
		table: blueprint.table.Clone(), actors: cloneActorsMap(blueprint.actors), nextToMove: nextToMove}
	return &c
}

func (state *RIGameState) dealPublicCard(card *Card) *RIGameState {

	c := createChild(state, state.round.NextRound(), DealPublicCardAction{card}, PlayerA, false)
	// important to deal using child deck / not current chance deck
	c.table.DropPublicCard(card)
	c.actors[ChanceId].(*Chance).deck.RemoveCard(card)
	return c
}

func (state *RIGameState) dealPrivateCards(cardA *Card, cardB *Card) *RIGameState {

	c := createChild(state, state.round.NextRound(), DealPrivateCardsAction{cardA, cardB}, PlayerA, false)
	// important to deal using child deck / not current chance deck
	c.playerActor(PlayerA).PlaceBet(c.table, Ante)
	c.actors[PlayerB].(*Player).PlaceBet(c.table, Ante)
	c.playerActor(PlayerA).CollectPrivateCard(cardA)
	c.actors[PlayerB].(*Player).CollectPrivateCard(cardB)
	c.actors[ChanceId].(*Chance).deck.RemoveCard(cardA)
	c.actors[ChanceId].(*Chance).deck.RemoveCard(cardB)

	return c
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
		player.Actions = []Action{}
		return player.Actions
	}

	bet := state.betSize()

	opponentStack := state.stack(player.Opponent())

	canBet := (player.Stack >= bet) && (opponentStack >= bet)
	canRaise := (player.Stack >= 2*bet) && (opponentStack >= 2*bet)

	// whenever betting roung is over (CALL OR CHECK->CHECK)
	bettingRoundEnded := state.causingAction.Name() == Call || (state.causingAction.Name() == Check && state.parent.causingAction.Name() == Check)
	if bettingRoundEnded {
		player.Actions = []Action{}
		return player.Actions
	}

	// single check implies BET or CHECK
	if state.causingAction.Name() == Check && state.parent.causingAction.Name() != Check {
		player.Actions = []Action{CheckAction}
		if canBet {
			player.Actions = append(player.Actions, BetAction)
		}
		return player.Actions
	}

	// RAISE/BET, you can CALL FOLD or RAISE (unless there has been 6 prior raises - 3 for each player)
	if state.causingAction.Name() == Bet || state.causingAction.Name() == Raise {
		player.Actions = []Action{CallAction, FoldAction}
		priorRaisesInCurrentRound := countPriorRaisesPerRound(state, state.round)
		if priorRaisesInCurrentRound < MaxRaises && canRaise {
			player.Actions = append(player.Actions, RaiseAction)
		}
		return player.Actions
	}

	if state.causingAction.Name() == DealPrivateCards || state.causingAction.Name() == DealPublicCards {
		player.Actions = []Action{CheckAction}
		if canBet {
			player.Actions = append(player.Actions, BetAction)
		}
		return player.Actions
	}
	panic(errors.New("code not reachable"))
}

func (state *RIGameState) playerActor(id ActorID) *Player {
	return state.actors[id].(*Player)
}

func (state *RIGameState) chanceActor() *Chance {
	return state.actors[ChanceId].(*Chance)
}
