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
		child = chance.dealPublicCard(state)
	}

	if move == DealPrivateCards {
		child = chance.dealPrivateCards(state)
	}
	return child
}

func (chance *Chance) dealPublicCard(state *RhodeIslandGameState) RhodeIslandGameState {
	child := state.CreateChild(state.round.NextRound(), DealPublicCard, state.table, PlayerA, false)
	child.table.DropPublicCard(chance.deck.DealNextCard())
	return child
}

func (chance *Chance) dealPrivateCards(state *RhodeIslandGameState) RhodeIslandGameState {
	table := state.table
	state.actors[PlayerA].(*PokerPlayer).PlaceBet(&table, Ante)
	state.actors[PlayerB].(*PokerPlayer).PlaceBet(&table, Ante)

	child := state.CreateChild(state.round.NextRound(), DealPrivateCards, table, PlayerA, false)
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
	betSize := state.BetSize()

	if move == Call || move == Bet || move == Raise {
		player.PlaceBet(&table, betSize)
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

func (player *PokerPlayer) PlaceBet(table *PokerTable, betSize float64) {
	table.AddToPot(betSize)
	player.stack -= betSize
}

func (player *PokerPlayer) GetAvailableMoves(state *RhodeIslandGameState) []Move {
	return player.availableMoves
}

func (player *PokerPlayer) computeAvailableActions(state *RhodeIslandGameState) {

	if player.availableMoves != nil {
		return
	}

	if state.causingMove == Fold {
		player.availableMoves = []Move{}
		return
	}
	betSize := state.BetSize()

	opponentStack := state.actors[player.Opponent()].(*PokerPlayer).stack

	allowedToBet := (player.stack >= betSize) && (opponentStack >= betSize)
	allowedToRaise := (player.stack >= 2*betSize) && (opponentStack >= 2*betSize)

	// whenever betting roung is over (CALL OR CHECK->CHECK)
	bettingRoundEnded := state.causingMove == Call || (state.causingMove == Check && state.parent.causingMove == Check)
	if bettingRoundEnded {
		player.availableMoves = nil
		return
	}

	// single check implies BET or CHECK
	if state.causingMove == Check && state.parent.causingMove != Check {
		player.availableMoves = []Move{Check}
		if allowedToBet {
			player.availableMoves = append(player.availableMoves, Bet)
		}
		return
	}

	// if RAISE/BET, you can CALL FOLD or RAISE (unless there has been 6 prior raises - 3 for each player)
	if state.causingMove == Bet || state.causingMove == Raise {
		player.availableMoves = []Move{Call, Fold}
		if countPriorRaises(*state) < 6 && allowedToRaise {
			player.availableMoves = append(player.availableMoves, Raise)
		}
		return
	}

	if state.causingMove == DealPrivateCards || state.causingMove == DealPublicCard {
		player.availableMoves = []Move{Check}
		if allowedToBet {
			player.availableMoves = append(player.availableMoves, Bet)
		}
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
