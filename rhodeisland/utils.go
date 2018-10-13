package rhodeisland

import . "github.com/int8/gopoker"

func cardsDiffersByTwo(cards []Card) bool {
	maxCard, minCard := int(C2), int(Ace)
	for _, card := range cards {
		if int(card.Name) >= maxCard {
			maxCard = int(card.Name)
		}

		if int(card.Name) <= minCard {
			minCard = int(card.Name)
		}
	}
	return maxCard-minCard == 2
}

func actionInSlice(a Action, actions []Action) bool {
	for _, x := range actions {
		if a == x {
			return true
		}
	}
	return false
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

func countPriorRaisesPerRound(node *RIGameState, round Round) int {
	if node == nil || node.causingAction.Name() != Raise || node.round != round {
		return 0
	}
	return 1 + countPriorRaisesPerRound(node.parent, round)
}
