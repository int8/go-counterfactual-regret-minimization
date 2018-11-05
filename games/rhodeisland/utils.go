package rhodeisland

import (
	"fmt"
	"github.com/int8/go-counterfactual-regret-minimization/acting"
	"github.com/int8/go-counterfactual-regret-minimization/cards"
	"github.com/int8/go-counterfactual-regret-minimization/games"
	"github.com/int8/go-counterfactual-regret-minimization/rounds"
)



func PrettyPrintInformationSet(infSet games.InformationSet) string {
	infSetArray := infSet.([InformationSetSizeBytes]byte)

	privateCardSymbol := cards.CardSymbol(read4BitsFromByteArray(infSetArray,0))
	privateCardSuit := cards.CardSuit(read3BitsFromByteArray(infSetArray,4))
	flopCardSymbol := cards.CardSymbol(read4BitsFromByteArray(infSetArray,7))
	flopCardSuit := cards.CardSuit(read3BitsFromByteArray(infSetArray,11))

	turnCardSymbol := cards.CardSymbol(read4BitsFromByteArray(infSetArray,14))
	turnCardSuit := cards.CardSuit(read3BitsFromByteArray(infSetArray,18))

	cardsString := fmt.Sprintf("%v%v %v%v %v%v ",privateCardSymbol, privateCardSuit, flopCardSymbol, flopCardSuit, turnCardSymbol, turnCardSuit)
	actionString := ""
	for i := 21; ; i += 3 {
		actionName := acting.ActionName(read3BitsFromByteArray(infSetArray, uint(i)))
		if actionName == acting.NoAction {
			break
		}
		actionString = fmt.Sprintf("%v ", actionName) + actionString
	}

	return cardsString + "| " + actionString
}

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



func read4BitsFromByteArray(data [InformationSetSizeBytes]byte, start uint) [4]bool {
	byteToRead := start / 8
	offset := start % 8
	result := [4]bool{}
	for i := 0; i < 4; i++ {
		if offset + uint(i) > 7 {
			result[i] = (data[byteToRead + 1] & (1 << ((offset + uint(i)) % 8))) > 0
		} else {
			result[i] = (data[byteToRead] & (1 << (offset + uint(i)))) > 0
		}
	}
	return result
}
func read3BitsFromByteArray(data [InformationSetSizeBytes]byte, start uint) [3]bool {
	byteToRead := start / 8
	offset := start % 8
	result := [3]bool{}
	for i := 0; i < 3; i++ {
		if offset + uint(i) > 7 {
			result[i] = (data[byteToRead + 1] & (1 << ((offset + uint(i)) % 8))) > 0
		} else {
			result[i] = (data[byteToRead] & (1 << (offset + uint(i))))  > 0
		}
	}
	return result
}