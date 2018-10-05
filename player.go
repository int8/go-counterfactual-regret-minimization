package gocfr

type PlayerIdentifier int8

const (
	PlayerA  PlayerIdentifier = 1
	PlayerB                   = -PlayerA
	ChanceId                  = 0
)

type ActionMaker interface {
	Act(state *RhodeIslandGameState, move Move) RhodeIslandGameState
	GetAvailableMoves(state *RhodeIslandGameState) []Move
}

type Chance struct {
	deck FullDeck
}

func (chance *Chance) Act(state *RhodeIslandGameState, move Move) RhodeIslandGameState {
	child := state.CreateChild(state.round.NextRound(), move, PlayerA, false)
	if move == DealPublicCard {
		child.table.DropPublicCard(chance.deck.DealNextCard())
	}

	if move == DealPrivateCards {
		child.actors[PlayerA].(*PokerPlayer).CollectPrivateCard(chance.deck.DealNextCard())
		child.actors[PlayerB].(*PokerPlayer).CollectPrivateCard(chance.deck.DealNextCard())
	}
	return child
}

func (chance *Chance) GetAvailableMoves(state *RhodeIslandGameState) []Move {
	if state.round == Start {
		return []Move{DealPrivateCards}
	}
	return []Move{DealPublicCard}
}

type PokerPlayer struct {
	id             PlayerIdentifier
	privateCards   []Card
	stack          float64
	availableMoves []Move
}

func (player *PokerPlayer) Act(state *RhodeIslandGameState, move Move) RhodeIslandGameState {
	if move == Fold {
		return state.CreateChild(state.round, move, player.Opponent(state).(*PokerPlayer).id, true)
	}
}

func (player *PokerPlayer) Opponent(state *RhodeIslandGameState) ActionMaker {
	return state.actors[-player.id]
}

func (player *PokerPlayer) CollectPrivateCard(card Card) {
	player.privateCards = append(player.privateCards, card)
}

func (player *PokerPlayer) GetAvailableMoves(state *RhodeIslandGameState) []Move {
	return player.availableMoves
}

func (player *PokerPlayer) computeAvailableActions(state *RhodeIslandGameState) {
	if player.availableMoves != nil {
		return
	}

	if state.causingMove == Fold {
		player.availableMoves = nil
		return
	}

	// whenever betting roung is over (CALL OR CHECK->CHECK)
	bettingRoundEnded := state.causingMove == Call || (state.causingMove == Check && state.parent.causingMove == Check)
	if bettingRoundEnded {
		player.availableMoves = nil
		return
	}

	// single check implies BET or CHECK
	if state.causingMove == Check && state.parent.causingMove != Check {
		player.availableMoves = []Move{Bet, Check}
		return
	}

	// you can only FOLD, RAISE or CALL on BET
	if state.causingMove == Bet {
		player.availableMoves = []Move{Call, Fold, Raise}
		return
	}

	// if RAISE, you can CALL FOLD or RAISE (unless there has been 6 prior raises - 3 for each player)
	if state.causingMove == Raise {
		if countPriorRaises(*state) < 6 {
			// allow raise if there has been less than 6 raises so far
			player.availableMoves = []Move{Call, Fold, Raise}
		} else {
			player.availableMoves = []Move{Call, Fold}
		}
		return
	}

	if state.causingMove == DealPrivateCards || state.causingMove == DealPublicCard {
		player.availableMoves = []Move{Bet, Check}
		return
	}
}

func (player PokerPlayer) String() string {
	if player.id == 1 {
		return "A"
	}
	if player.id == -1 {
		return "B"
	}
	return "Chance"
}
