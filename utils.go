package gocfr

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func selectActionByMove(actions []Action, move Move) int {
	for i, a := range actions {
		if a.move == move {
			return i
		}
	}
	return -1
}

func countPriorRaises(node RhodeIslandGameState) int {
	if &node == nil || node.causingMove != Raise {
		return 0
	}
	return 1 + countPriorRaises(*node.parent)

}

func roundCheckFunc(expectedRound Round) func(node RhodeIslandGameState) bool {
	return func(node RhodeIslandGameState) bool { return node.round == expectedRound }
}

func GameEndFunc() func(state RhodeIslandGameState) bool {
	return func(state RhodeIslandGameState) bool { return state.IsTerminal() }
}

func NoRaiseAvailable() func(state RhodeIslandGameState) bool {
	return func(state RhodeIslandGameState) bool {
		moves := state.CurrentActor().GetAvailableMoves(&state)
		for _, m := range moves {
			if m == Raise {
				return false
			}
		}
		return true
	}
}

func ActionMakerToMoveFunc(actionMakerId ActionMakerIdentifier) func(state RhodeIslandGameState) bool {
	return func(state RhodeIslandGameState) bool {
		return state.nextToMove == actionMakerId
	}
}
