package kuhn

import . "github.com/int8/gopoker"

func actionInSlice(a Action, actions []Action) bool {
	for _, x := range actions {
		if a == x {
			return true
		}
	}
	return false
}

func cloneActorsMap(srcActors map[ActorID]Actor) map[ActorID]Actor {
	actors := make(map[ActorID]Actor)
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
