package gocfr

type ActionMakerIdentifier int8

const (
	PlayerA       ActionMakerIdentifier = 1
	PlayerB                             = -PlayerA
	ChanceId                            = 0
	NoActionMaker                       = 100
)

type ActionMaker interface {
	Act(state *RhodeIslandGameState, move Move) RhodeIslandGameState
	GetAvailableMoves(state *RhodeIslandGameState) []Move
}

type Chance struct {
	id   ActionMakerIdentifier
	deck FullDeck
}

func (chance *Chance) Act(state *RhodeIslandGameState, move Move) (child RhodeIslandGameState) {

	if move == DealPublicCard {
		child = chance.dealPublicCard(state, move)
	}

	if move == DealPrivateCards {
		child = chance.dealPrivateCards(state, move)
	}
	return child
}

func (chance *Chance) dealPublicCard(state *RhodeIslandGameState, move Move) RhodeIslandGameState {
	child := state.CreateChild(state.round.NextRound(), move, state.table, PlayerA, false)
	child.table.DropPublicCard(chance.deck.DealNextCard())
	return child
}

func (chance *Chance) dealPrivateCards(state *RhodeIslandGameState, move Move) RhodeIslandGameState {
	table := state.table
	state.actors[PlayerA].(*PokerPlayer).stack -= Ante
	state.actors[PlayerB].(*PokerPlayer).stack -= Ante
	table.potSize += 2 * Ante
	child := state.CreateChild(state.round.NextRound(), move, table, PlayerA, false)
	child.actors[PlayerA].(*PokerPlayer).CollectPrivateCard(chance.deck.DealNextCard())
	child.actors[PlayerB].(*PokerPlayer).CollectPrivateCard(chance.deck.DealNextCard())
	return child
}

func (chance *Chance) GetAvailableMoves(state *RhodeIslandGameState) []Move {
	if state.round == Start {
		return []Move{DealPrivateCards}
	}
	if !state.terminal {
		return []Move{DealPublicCard}
	}
	return nil
}

type PokerPlayer struct {
	id             ActionMakerIdentifier
	privateCards   []Card
	stack          float64
	availableMoves []Move
}

func (player *PokerPlayer) Act(state *RhodeIslandGameState, move Move) RhodeIslandGameState {

	table := state.table
	var betSize float64

	if state.round < Flop {
		betSize = PreFlopBetSize
	} else {
		betSize = PostFlopBetSize
	}

	if move == Call || move == Bet {
		player.stack -= betSize
		table.potSize += betSize
	}

	if move == Raise {
		player.stack -= 2 * betSize
		table.potSize += 2 * betSize
	}

	if move == Fold || (state.round == Turn && (move == Call || (move == Check && state.causingMove == Check))) {
		return state.CreateChild(state.round, move, table, NoActionMaker, true)
	}

	if move == Call || (move == Check && state.causingMove == Check) {
		return state.CreateChild(state.round, move, table, ChanceId, false)
	}

	return state.CreateChild(state.round, move, table, player.Opponent(), false)
}

func (player *PokerPlayer) Opponent() ActionMakerIdentifier {
	return -player.id
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

	// if RAISE/BET, you can CALL FOLD or RAISE (unless there has been 6 prior raises - 3 for each player)
	if state.causingMove == Bet || state.causingMove == Raise {
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
