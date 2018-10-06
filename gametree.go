package gocfr

// GameState
type GameState struct {
	round       Round
	parent      *GameState
	causingMove Move
	table       *Table
	actors      map[ActorId]Actor
	nextToMove  ActorId
	dealer      ActorId
	terminal    bool
}

func (state *GameState) CurrentActor() Actor {
	return state.actors[state.nextToMove]
}

func (state *GameState) BetSize() float64 {
	if state.round < Flop {
		return PreFlopBetSize
	}
	return PostFlopBetSize
}

func CreateRoot(dealer ActorId, playerAStack float64, playerBStack float64) *GameState {
	playerA := &Player{id: PlayerA, moves: nil, cards: []Card{}, stack: playerAStack}
	playerB := &Player{id: PlayerB, moves: nil, cards: []Card{}, stack: playerBStack}
	chance := &Chance{id: ChanceId, deck: CreateFullDeck(true)}

	actors := map[ActorId]Actor{PlayerA: playerA, PlayerB: playerB, ChanceId: chance}
	table := &Table{pot: 0, cards: []Card{}}

	return &GameState{round: Start, table: table,
		actors: actors, nextToMove: ChanceId, causingMove: NoMove, dealer: dealer}
}

func (state *GameState) CreateChild(round Round, move Move, nextToMove ActorId, terminal bool) *GameState {
	child := GameState{round: round,
		parent: state, causingMove: move, terminal: terminal,
		table: state.table.Clone(), actors: cloneActorsMap(state.actors), nextToMove: nextToMove, dealer: state.dealer}
	return &child
}

func (state *GameState) IsTerminal() bool {
	return state.terminal
}

func createRootWithPreparedDeck(dealer ActorId, playerAStack float64, playerBStack float64, deck *FullDeck) *GameState {
	playerA := &Player{id: PlayerA, moves: nil, cards: []Card{}, stack: playerAStack}
	playerB := &Player{id: PlayerB, moves: nil, cards: []Card{}, stack: playerBStack}
	chance := &Chance{id: ChanceId, deck: deck}

	actors := map[ActorId]Actor{PlayerA: playerA, PlayerB: playerB, ChanceId: chance}
	table := &Table{pot: 0, cards: []Card{}}

	return &GameState{round: Start, table: table,
		actors: actors, nextToMove: ChanceId, causingMove: NoMove, dealer: dealer}
}
