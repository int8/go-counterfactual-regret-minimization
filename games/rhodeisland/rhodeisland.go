package rhodeisland

import (
	"errors"
	"github.com/int8/go-counterfactual-regret-minimization/acting"
	"github.com/int8/go-counterfactual-regret-minimization/cards"
	"github.com/int8/go-counterfactual-regret-minimization/games"
	"github.com/int8/go-counterfactual-regret-minimization/rounds"
	"github.com/int8/go-counterfactual-regret-minimization/table"
)

const PreFlopBetSize float32 = 10.
const PostFlopBetSize float32 = 20.

const InformationSetSize = 8 * 12
const InformationSetSizeBytes = 12

var MaxRaises = 3

const Ante float32 = 5.0

// RIGameState - RhodeIslandGameState
type RIGameState struct {
	round         rounds.PokerRound
	parent        *RIGameState
	causingAction acting.Action
	table         *table.PokerTable
	actors        map[acting.ActorID]acting.Actor
	nextToMove    acting.ActorID
	terminal      bool
}

func (state *RIGameState) Act(action acting.Action) games.GameState {
	switch state.actors[state.nextToMove].(type) {
	case *Chance:
		return state.actAsChance(action)
	case *Player:
		return state.actAsPlayer(action)
	}
	return nil
}

func (state *RIGameState) Actions() []acting.Action {

	switch state.actors[state.nextToMove].(type) {
	case *Chance:
		return state.chanceActions(state.chanceActor())
	case *Player:
		return state.playerActions(state.playerActor(state.nextToMove))
	}
	return nil
}

func (state *RIGameState) IsChance() bool {
	return state.nextToMove == acting.ChanceId
}

func (state *RIGameState) IsTerminal() bool {
	return state.terminal
}

func (state *RIGameState) Parent() games.GameState {
	return state.parent
}

func (state *RIGameState) CurrentActor() acting.Actor {
	return state.actors[state.nextToMove]
}

func (state *RIGameState) Evaluate() float32 {
	actor := state.playerActor(state.CurrentActor().GetID())
	opponent := state.playerActor(-state.CurrentActor().GetID())
	if state.IsTerminal() {
		if state.causingAction.Name() == acting.Fold {
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

func (state *RIGameState) InformationSet() games.InformationSet {

	privateCard := cards.Card{Symbol: state.playerActor(state.nextToMove).Card.Symbol, Suit: state.playerActor(state.nextToMove).Card.Suit}
	flopCard, turnCard := cards.NoCard, cards.NoCard

	if len(state.table.Cards) > 0 {
		flopCard = cards.Card{Symbol: state.table.Cards[0].Symbol, Suit: state.table.Cards[0].Suit}
	}

	if len(state.table.Cards) > 1 {
		turnCard = cards.Card{Symbol: state.table.Cards[1].Symbol, Suit: state.table.Cards[1].Suit}
	}
	informationSet := [InformationSetSizeBytes]byte{}

	infSetBool := [InformationSetSize]bool{
		privateCard.Symbol[0], privateCard.Symbol[1], privateCard.Symbol[2], privateCard.Symbol[3],
		privateCard.Suit[0], privateCard.Suit[1], privateCard.Suit[2],
		flopCard.Symbol[0], flopCard.Symbol[1], flopCard.Symbol[2], flopCard.Symbol[3],
		flopCard.Suit[0], flopCard.Suit[1], flopCard.Suit[2],
		turnCard.Symbol[0], turnCard.Symbol[1], turnCard.Symbol[2], turnCard.Symbol[3],
		turnCard.Suit[0], turnCard.Suit[1], turnCard.Suit[2],
	}

	currentState := state
	for i := 21; currentState.round != rounds.Start; i += 3 {
		actionName := currentState.causingAction.Name()
		infSetBool[i] = actionName[0]
		infSetBool[i+1] = actionName[1]
		infSetBool[i+2] = actionName[2]

		currentState = currentState.parent
		if currentState == nil {
			break
		}
	}

	for i := 0; i < InformationSetSizeBytes; i++ {
		informationSet[i] = acting.CreateByte(infSetBool[(i * 8):((i + 1) * 8)])
	}

	return games.InformationSet(informationSet)
}

func (state *RIGameState) stack(id acting.ActorID) float32 {
	return state.actors[id].(*Player).Stack
}

func (state *RIGameState) actAsChance(action acting.Action) games.GameState {
	var c *RIGameState
	if action.Name() == acting.DealPublicCards {
		c = state.dealPublicCard(action.(DealPublicCardAction).Card)
	}

	if action.Name() == acting.DealPrivateCards {
		c = state.dealPrivateCards(action.(DealPrivateCardsAction).CardA, action.(DealPrivateCardsAction).CardB)
	}
	return c
}

func (state *RIGameState) actAsPlayer(action acting.Action) games.GameState {

	var c *RIGameState

	if !actionInSlice(action, state.Actions()) {
		panic("action not available")
	}
	actor := state.playerActor(state.nextToMove)
	betSize := state.betSize()

	defer func() {
		if action.Name() == acting.Call || action.Name() == acting.Bet {
			c.playerActor(actor.GetID()).PlaceBet(c.table, betSize)
		}
		if action.Name() == acting.Raise {
			c.playerActor(actor.GetID()).PlaceBet(c.table, 2*betSize)
		}
		if action.Name() == acting.Fold {
			c.playerActor(-actor.GetID()).PlaceBet(c.table, -betSize)
		}
	}()

	if action.Name() == acting.Fold || (state.round == rounds.Turn && (action.Name() == acting.Call || (action.Name() == acting.Check && state.causingAction.Name() == acting.Check))) {
		c = createChild(state, state.round, action, actor.Opponent(), true)
		return c
	}

	if action.Name() == acting.Call || (action.Name() == acting.Check && state.causingAction.Name() == acting.Check) {
		c = createChild(state, state.round, action, acting.ChanceId, false)
		return c
	}

	c = createChild(state, state.round, action, actor.Opponent(), false)
	return c

}

func (state *RIGameState) betSize() float32 {
	if state.round < rounds.Flop {
		return PreFlopBetSize
	}
	return PostFlopBetSize
}

func Root(playerA *Player, playerB *Player, deck cards.Deck) *RIGameState {
	chance := &Chance{id: acting.ChanceId, deck: deck}

	actors := map[acting.ActorID]acting.Actor{acting.PlayerA: playerA, acting.PlayerB: playerB, acting.ChanceId: chance}
	pokerTable := &table.PokerTable{Pot: 0, Cards: []cards.Card{}}
	return &RIGameState{round: rounds.Start, table: pokerTable,
		actors: actors, nextToMove: acting.ChanceId, causingAction: nil}
}

func createChild(blueprint *RIGameState, round rounds.PokerRound, action acting.Action, nextToMove acting.ActorID, terminal bool) *RIGameState {
	c := RIGameState{round: round,
		parent: blueprint, causingAction: action, terminal: terminal,
		table: blueprint.table.Clone(), actors: cloneActorsMap(blueprint.actors), nextToMove: nextToMove}
	return &c
}

func (state *RIGameState) dealPublicCard(card *cards.Card) *RIGameState {

	c := createChild(state, state.round.NextRound(), DealPublicCardAction{card}, acting.PlayerA, false)
	// important to deal using child deck / not current chance deck
	c.table.DropPublicCard(card)
	c.actors[acting.ChanceId].(*Chance).deck.RemoveCard(card)
	return c
}

func (state *RIGameState) dealPrivateCards(cardA *cards.Card, cardB *cards.Card) *RIGameState {

	c := createChild(state, state.round.NextRound(), DealPrivateCardsAction{cardA, cardB}, acting.PlayerA, false)
	// important to deal using child deck / not current chance deck
	c.playerActor(acting.PlayerA).PlaceBet(c.table, Ante)
	c.actors[acting.PlayerB].(*Player).PlaceBet(c.table, Ante)
	c.playerActor(acting.PlayerA).CollectPrivateCard(cardA)
	c.actors[acting.PlayerB].(*Player).CollectPrivateCard(cardB)
	c.actors[acting.ChanceId].(*Chance).deck.RemoveCard(cardA)
	c.actors[acting.ChanceId].(*Chance).deck.RemoveCard(cardB)

	return c
}

func (state *RIGameState) chanceActions(chance *Chance) []acting.Action {
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

	actions := make([]acting.Action, chance.deck.CardsLeft())
	remainingCards := chance.deck.RemainingCards()
	for i, card := range remainingCards {
		actions[i] = DealPublicCardAction{card}
	}
	return actions
}

func (state *RIGameState) playerActions(player *Player) []acting.Action {

	if state.causingAction.Name() == acting.Fold {
		player.Actions = []acting.Action{}
		return player.Actions
	}

	bet := state.betSize()

	opponentStack := state.stack(player.Opponent())

	canBet := (player.Stack >= bet) && (opponentStack >= bet)
	canRaise := (player.Stack >= 2*bet) && (opponentStack >= 2*bet)

	// whenever betting roung is over (CALL OR CHECK->CHECK)
	bettingRoundEnded := state.causingAction.Name() == acting.Call || (state.causingAction.Name() == acting.Check && state.parent.causingAction.Name() == acting.Check)
	if bettingRoundEnded {
		player.Actions = []acting.Action{}
		return player.Actions
	}

	// single check implies BET or CHECK
	if state.causingAction.Name() == acting.Check && state.parent.causingAction.Name() != acting.Check {
		player.Actions = []acting.Action{CheckAction}
		if canBet {
			player.Actions = append(player.Actions, BetAction)
		}
		return player.Actions
	}

	// RAISE/BET, you can CALL FOLD or RAISE (unless there has been 6 prior raises - 3 for each player)
	if state.causingAction.Name() == acting.Bet || state.causingAction.Name() == acting.Raise {
		player.Actions = []acting.Action{CallAction, FoldAction}
		priorRaisesInCurrentRound := countPriorRaisesPerRound(state, state.round)
		if priorRaisesInCurrentRound < MaxRaises && canRaise {
			player.Actions = append(player.Actions, RaiseAction)
		}
		return player.Actions
	}

	if state.causingAction.Name() == acting.DealPrivateCards || state.causingAction.Name() == acting.DealPublicCards {
		player.Actions = []acting.Action{CheckAction}
		if canBet {
			player.Actions = append(player.Actions, BetAction)
		}
		return player.Actions
	}
	panic(errors.New("code not reachable"))
}

func (state *RIGameState) playerActor(id acting.ActorID) *Player {
	return state.actors[id].(*Player)
}

func (state *RIGameState) chanceActor() *Chance {
	return state.actors[acting.ChanceId].(*Chance)
}
