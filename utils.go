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

func cardsDiffersByTwo(cards []Card) bool {
	maxCard, minCard := int(C2), int(Ace)
	for _, card := range cards {
		if int(card.name) >= maxCard {
			maxCard = int(card.name)
		}

		if int(card.name) <= minCard {
			minCard = int(card.name)
		}
	}
	return maxCard-minCard == 2
}
