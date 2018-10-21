package rhodeisland

import (
	"fmt"
	. "github.com/int8/gopoker"
)

func PrettyPrintInformationSet(infSet InformationSet) string {
	infSetArray := infSet.([InformationSetSize]bool)
	prvCardName := CardName([4]bool{infSetArray[0], infSetArray[1], infSetArray[2], infSetArray[3]})
	prvCardColor := CardSuit([3]bool{infSetArray[4], infSetArray[5], infSetArray[6]})
	flopCardName := CardName([4]bool{infSetArray[7], infSetArray[8], infSetArray[9], infSetArray[10]})
	flopCardColor := CardSuit([3]bool{infSetArray[11], infSetArray[12], infSetArray[13]})
	turnCardName := CardName([4]bool{infSetArray[14], infSetArray[15], infSetArray[16], infSetArray[17]})
	turnCardColor := CardSuit([3]bool{infSetArray[18], infSetArray[19], infSetArray[20]})

	cardsString := fmt.Sprintf("%v%v* %v%v%v%v", prvCardName, prvCardColor, flopCardName, flopCardColor, turnCardName, turnCardColor)

	actionString := ""
	for i := 21; ; i += 3 {
		actionName := ActionName([3]bool{infSetArray[i], infSetArray[i+1], infSetArray[i+2]})
		if actionName == NoAction {
			break
		}
		actionString = fmt.Sprintf("%v ", actionName) + actionString
	}

	return cardsString + "| " + actionString
}

func cardsDiffersByTwo(cards []Card) bool {
	maxCard, minCard := CardNameInt(C2), CardNameInt(Ace)
	for _, card := range cards {
		cardInt := CardNameInt(card.Name)
		if cardInt >= maxCard {
			maxCard = cardInt
		}

		if cardInt <= minCard {
			minCard = cardInt
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
