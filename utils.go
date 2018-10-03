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
	if &node == nil || node.causingAction.move != Raise {
		return 0
	} else {
		return 1 + countPriorRaises(*node.parent)
	}
}
