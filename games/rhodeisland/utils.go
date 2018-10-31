package rhodeisland

import (
	"github.com/int8/go-counterfactual-regret-minimization/acting"
	"github.com/int8/go-counterfactual-regret-minimization/cards"
	"github.com/int8/go-counterfactual-regret-minimization/rounds"
)

//
//func PrettyPrintInformationSet(infSet games.InformationSet) string {
//	infSetArray := infSet.([InformationSetSizeBytes]byte)
//	prvCardSymbol := cards.CardSymbol([4]bool{infSetArray[0], infSetArray[1], infSetArray[2], infSetArray[3]})
//	prvCardColor := cards.CardSuit([3]bool{infSetArray[4], infSetArray[5], infSetArray[6]})
//	flopCardSymbol := cards.CardSymbol([4]bool{infSetArray[7], infSetArray[8], infSetArray[9], infSetArray[10]})
//	flopCardColor := cards.CardSuit([3]bool{infSetArray[11], infSetArray[12], infSetArray[13]})
//	turnCardSymbol := cards.CardSymbol([4]bool{infSetArray[14], infSetArray[15], infSetArray[16], infSetArray[17]})
//	turnCardColor := cards.CardSuit([3]bool{infSetArray[18], infSetArray[19], infSetArray[20]})
//
//	cardsString := fmt.Sprintf("%v%v* %v%v%v%v", prvCardSymbol, prvCardColor, flopCardSymbol, flopCardColor, turnCardSymbol, turnCardColor)
//
//	actionString := ""
//	for i := 21; ; i += 3 {
//		actionName := acting.ActionName([3]bool{infSetArray[i], infSetArray[i+1], infSetArray[i+2]})
//		if actionName == acting.NoAction {
//			break
//		}
//		actionString = fmt.Sprintf("%v ", actionName) + actionString
//	}
//
//	return cardsString + "| " + actionString
//}

func cardsDiffersByTwo(inputCards []cards.Card) bool {
	maxCard, minCard := cards.CardSymbol2Int(cards.C2), cards.CardSymbol2Int(cards.Ace)
	for _, card := range inputCards {
		cardInt := cards.CardSymbol2Int(card.Symbol)
		if cardInt >= maxCard {
			maxCard = cardInt
		}

		if cardInt <= minCard {
			minCard = cardInt
		}
	}
	return maxCard-minCard == 2
}

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

func countPriorRaisesPerRound(node *RIGameState, round rounds.PokerRound) int {
	if node == nil || node.causingAction.Name() != acting.Raise || node.round != round {
		return 0
	}
	return 1 + countPriorRaisesPerRound(node.parent, round)
}
