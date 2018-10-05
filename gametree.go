package gocfr

// RhodeIslandGameState
type RhodeIslandGameState struct {
	round       Round
	parent      *RhodeIslandGameState
	causingMove Move
	table       PokerTable
	actors      map[PlayerIdentifier]ActionMaker
	nextToMove  PlayerIdentifier
	terminal    bool
}

func CreateRoot(playerAStack float64, playerBStack float64) RhodeIslandGameState {

	actors := map[PlayerIdentifier]ActionMaker{PlayerA: &PokerPlayer{}, PlayerB: &PokerPlayer{}, ChanceId: &Chance{}}
	table := PokerTable{potSize: 0, publicCards: []Card{}}

	return RhodeIslandGameState{round: Start, table: table,
		actors: actors, nextToMove: ChanceId, causingMove: NoMove}
}

func (node *RhodeIslandGameState) CreateChild(round Round, move Move, nextToMove PlayerIdentifier, terminal bool) RhodeIslandGameState {
	child := RhodeIslandGameState{round: round,
		parent: node, causingMove: move, terminal: terminal,
		table: node.table, actors: node.actors, nextToMove: nextToMove}
	return child
}

func (node *RhodeIslandGameState) IsTerminal() bool {
	actions := node.actors[node.nextToMove].GetAvailableMoves(node)
	return actions == nil
}
