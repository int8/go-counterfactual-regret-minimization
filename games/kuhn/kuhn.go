package kuhn

import (
	"errors"
	"github.com/int8/gopoker/acting"
	"github.com/int8/gopoker/cards"
	"github.com/int8/gopoker/games"
	"github.com/int8/gopoker/rounds"
	"github.com/int8/gopoker/table"
)

const BetSize float32 = 1.0
const Ante float32 = 1.0
const InformationSetSize = 24
const InformationSetSizeBytes = 3

// KuhnGameState - Kuhn Poker Game State
type KuhnGameState struct {
	round         rounds.PokerRound
	parent        *KuhnGameState
	causingAction acting.Action
	table         *table.PokerTable
	actors        map[acting.ActorID]acting.Actor
	nextToMove    acting.ActorID
	terminal      bool
}

func (state *KuhnGameState) Act(action acting.Action) games.GameState {
	switch state.actors[state.nextToMove].(type) {
	case *Chance:
		return state.actAsChance(action)
	case *Player:
		return state.actAsPlayer(action)
	}
	return nil
}

func (state *KuhnGameState) Actions() []acting.Action {

	switch state.actors[state.nextToMove].(type) {
	case *Chance:
		return state.chanceActions(state.actors[state.nextToMove].(*Chance))
	case *Player:
		return state.playerActions(state.actors[state.nextToMove].(*Player))
	}
	return nil
}

func (state *KuhnGameState) IsChance() bool {
	return state.nextToMove == acting.ChanceId
}

func (state *KuhnGameState) IsTerminal() bool {
	return state.terminal
}

func (state *KuhnGameState) Parent() games.GameState {
	return state.parent
}

func (state *KuhnGameState) CurrentActor() acting.Actor {
	return state.actors[state.nextToMove]
}

//TODO: test it carefully
func (state *KuhnGameState) Evaluate() float32 {
	currentActor := state.playerActor(state.CurrentActor().GetID())
	currentActorOpponent := state.playerActor(-state.CurrentActor().GetID())
	if state.IsTerminal() {
		if state.causingAction.Name() == acting.Fold {
			currentActor.UpdateStack(currentActor.Stack + state.table.Pot)
			return float32(currentActor.GetID()) * (state.table.Pot / 2)
		}
		currentActorHandValue := currentActor.EvaluateHand(state.table)
		currentActorOpponentHandValue := currentActorOpponent.EvaluateHand(state.table)

		if currentActorHandValue > currentActorOpponentHandValue {
			currentActor.UpdateStack(currentActor.Stack + state.table.Pot)
			return float32(currentActor.GetID()) * (state.table.Pot / 2)
		} else {
			currentActorOpponent.UpdateStack(currentActorOpponent.Stack + state.table.Pot)
			return float32(currentActorOpponent.GetID()) * (state.table.Pot / 2)
		}
	}
	currentActor.UpdateStack(currentActor.Stack + state.table.Pot/2)
	currentActorOpponent.UpdateStack(currentActorOpponent.Stack + state.table.Pot/2)
	return 0.0
}

func (state *KuhnGameState) InformationSet() games.InformationSet {

	privateCardSymbol := state.actors[state.nextToMove].(*Player).Card.Symbol
	privateCardSuit := state.actors[state.nextToMove].(*Player).Card.Suit
	informationSet := [InformationSetSizeBytes]byte{}

	informationSetBool := [InformationSetSize]bool{
		privateCardSymbol[0], privateCardSymbol[1], privateCardSymbol[2], privateCardSymbol[3],
		privateCardSuit[0], privateCardSuit[1], privateCardSuit[2],
	}
	currentState := state
	for i := 7; currentState.round != rounds.Start; i += 3 {
		actionName := currentState.causingAction.Name()
		informationSetBool[i] = actionName[0]
		informationSetBool[i+1] = actionName[1]
		informationSetBool[i+2] = actionName[2]
		currentState = currentState.parent
		if currentState == nil {
			break
		}
	}
	for i := 0; i < InformationSetSizeBytes; i++ {
		informationSet[i] = acting.CreateByte(informationSetBool[(i * 8):((i + 1) * 8)])
	}

	return games.InformationSet(informationSet)
}

func (state *KuhnGameState) stack(actor acting.ActorID) float32 {
	return state.actors[actor].(*Player).Stack
}

func (state *KuhnGameState) actAsChance(action acting.Action) games.GameState {
	var child *KuhnGameState
	if action.Name() == acting.DealPrivateCards {
		child = state.dealPrivateCards(action.(DealPrivateCardsAction).CardA, action.(DealPrivateCardsAction).CardB)
	}
	return child
}

func (state *KuhnGameState) actAsPlayer(action acting.Action) games.GameState {

	var child *KuhnGameState

	if !actionInSlice(action, state.Actions()) {
		panic("action not available")
	}
	actor := state.CurrentActor()

	defer func() {
		if action.Name() == acting.Call || action.Name() == acting.Bet {
			child.playerActor(actor.GetID()).PlaceBet(child.table, BetSize)
		}
		if action.Name() == acting.Fold {
			//opponent of folding player can now take his bet back
			child.actors[state.CurrentActor().(*Player).Opponent()].(*Player).PlaceBet(child.table, -BetSize)
		}
	}()

	if action.Name() == acting.Fold || action.Name() == acting.Call || (action.Name() == acting.Check && state.causingAction.Name() == acting.Check) {
		child = createChild(state, state.round, action, state.CurrentActor().(*Player).Opponent(), true)
		return child
	}

	child = createChild(state, state.round, action, state.CurrentActor().(*Player).Opponent(), false)
	return child
}

func (state *KuhnGameState) betSize() float32 {
	return 1.0
}

func Root(playerA *Player, playerB *Player) *KuhnGameState {
	chance := &Chance{id: acting.ChanceId, deck: CreateKuhnDeck()}

	actors := map[acting.ActorID]acting.Actor{acting.PlayerA: playerA, acting.PlayerB: playerB, acting.ChanceId: chance}
	table := &table.PokerTable{Pot: 0, Cards: []cards.Card{}}

	return &KuhnGameState{round: rounds.Start, table: table,
		actors: actors, nextToMove: acting.ChanceId, causingAction: nil}
}

func createChild(blueprint *KuhnGameState, round rounds.PokerRound, Action acting.Action, nextToMove acting.ActorID, terminal bool) *KuhnGameState {
	child := KuhnGameState{round: round,
		parent: blueprint, causingAction: Action, terminal: terminal,
		table: blueprint.table.Clone(), actors: cloneActorsMap(blueprint.actors), nextToMove: nextToMove}
	return &child
}

func (state *KuhnGameState) dealPrivateCards(cardA *cards.Card, cardB *cards.Card) *KuhnGameState {

	child := createChild(state, state.round.NextRound(), DealPrivateCardsAction{cardA, cardB}, acting.PlayerA, false)
	// important to deal using child deck / not current chance deck
	child.actors[acting.PlayerA].(*Player).PlaceBet(child.table, Ante)
	child.actors[acting.PlayerB].(*Player).PlaceBet(child.table, Ante)
	child.actors[acting.PlayerA].(*Player).CollectPrivateCard(cardA)
	child.actors[acting.PlayerB].(*Player).CollectPrivateCard(cardB)
	child.actors[acting.ChanceId].(*Chance).deck.RemoveCard(cardA)
	child.actors[acting.ChanceId].(*Chance).deck.RemoveCard(cardB)

	return child
}

func (state *KuhnGameState) chanceActions(chance *Chance) []acting.Action {
	if state.round == rounds.Start {
		deckSize := int(chance.deck.CardsLeft())
		actions := make([]acting.Action, deckSize*(deckSize-1))
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
	return nil
}

func (state *KuhnGameState) playerActions(player *Player) []acting.Action {

	if state.causingAction.Name() == acting.Fold {
		player.Actions = []acting.Action{}
		return player.Actions
	}

	betSize := state.betSize()
	opponentStack := state.stack(player.Opponent())
	allowedToBet := (player.Stack >= betSize) && (opponentStack >= betSize)

	// whenever betting round is over (CALL OR CHECK->CHECK)
	bettingRoundEnded := state.causingAction.Name() == acting.Call || (state.causingAction.Name() == acting.Check && state.parent.causingAction.Name() == acting.Check)
	if bettingRoundEnded {
		player.Actions = []acting.Action{}
		return player.Actions
	}

	// single check implies BET or CHECK
	if state.causingAction.Name() == acting.Check && state.parent.causingAction.Name() != acting.Check {
		player.Actions = []acting.Action{CheckAction}
		if allowedToBet {
			player.Actions = append(player.Actions, BetAction)
		}
		return player.Actions
	}

	if state.causingAction.Name() == acting.Bet {
		player.Actions = []acting.Action{CallAction, FoldAction}
		return player.Actions
	}

	if state.causingAction.Name() == acting.DealPrivateCards || state.causingAction.Name() == acting.DealPublicCards {
		player.Actions = []acting.Action{CheckAction}
		if allowedToBet {
			player.Actions = append(player.Actions, BetAction)
		}
		return player.Actions
	}
	panic(errors.New("Code not reachable."))
}

func (state *KuhnGameState) playerActor(id acting.ActorID) *Player {
	return state.actors[id].(*Player)
}
