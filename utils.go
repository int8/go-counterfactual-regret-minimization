package gocfr

func makeRange(min, max uint8) []uint8 {
	a := make([]uint8, max-min+1)
	for i := range a {
		a[i] = min + uint8(i)
	}
	return a
}

func cloneActorsMap(srcActors map[ActorId]Actor) map[ActorId]Actor {
	actors := make(map[ActorId]Actor)
	for id, actor := range srcActors {
		switch actor.(type) {
		case *Player:
			actors[id] = actor.(*Player).Clone()
		case *Chance:
			actors[id] = actor.(*Chance).Clone()
		}
	}
	return actors
}

func countPriorRaises(node *GameState) int {
	if node == nil || node.causingMove != Raise {
		return 0
	}
	return 1 + countPriorRaises(node.parent)

}

func roundCheck(expectedRound Round) func(node *GameState) bool {
	return func(node *GameState) bool { return node.round == expectedRound }
}

func gameEnd() func(state *GameState) bool {
	return func(state *GameState) bool { return state.IsTerminal() }
}

func noRaiseAvailable() func(state *GameState) bool {
	return func(state *GameState) bool {
		moves := state.CurrentActor().GetAvailableMoves(state)
		for _, m := range moves {
			if m == Raise {
				return false
			}
		}
		return true
	}
}

func actorToMove(actorId ActorId) func(state *GameState) bool {
	return func(state *GameState) bool {
		return state.nextToMove == actorId
	}
}

func stackEqualTo(player ActorId, stack float64) func(state *GameState) bool {
	return func(state *GameState) bool {
		return state.actors[player].(*Player).stack == stack
	}
}

func noTest() func(state *GameState) bool {
	return func(state *GameState) bool {
		return true
	}
}

func onlyCheckAvailable() func(state *GameState) bool {
	return func(state *GameState) bool {
		moves := state.CurrentActor().GetAvailableMoves(state)
		if len(moves) == 1 && moves[0] == Check {
			return true
		}
		return false
	}
}

func checkAndBetAvailable() func(state *GameState) bool {
	return func(state *GameState) bool {
		moves := state.CurrentActor().GetAvailableMoves(state)
		if len(moves) == 2 && moves[0] == Check && moves[1] == Bet {
			return true
		}
		return false
	}
}
