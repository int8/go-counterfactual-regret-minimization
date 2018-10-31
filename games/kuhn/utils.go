package kuhn

import (
	"github.com/int8/gopoker/acting"
)

func actionInSlice(a acting.Action, actions []acting.Action) bool {
	for _, x := range actions {
		if a == x {
			return true
		}
	}
	return false
}

func cloneActorsMap(srcActors map[acting.ActorID]acting.Actor) map[acting.ActorID]acting.Actor {
	actors := make(map[acting.ActorID]acting.Actor)
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
