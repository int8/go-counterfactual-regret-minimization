package gocfr

import (
	"reflect"
	"testing"
)

func TestFullDeckCardsCount(t *testing.T) {
	fullDeck := CreateFullDeck()

	if len(fullDeck.Cards) != 52 {
		t.Error("Full deck should count 52 cards")
	}

	if fullDeck.CardsLeft() != 52 {
		t.Errorf("Full deck should have 52 cards left after initialization but have %v", fullDeck.CardsLeft())
	}

	for i := range fullDeck.Cards {
		if fullDeck.CardsLeft() != 52-i {
			t.Errorf("Full deck should have %v cards left dealing %v card but have %v", 52-i, i, fullDeck.CardsLeft())
		}
		fullDeck.DealNextCard()
	}

}

func TestFullDeckCardsShuffling(t *testing.T) {
	fullDeck := CreateFullDeck()
	orderBeforeShuffling := make([]int, 52, 52)
	copy(orderBeforeShuffling, fullDeck.shuffleOrder)

	for range fullDeck.Cards {
		fullDeck.DealNextCard()
	}

	if reflect.DeepEqual(fullDeck.shuffleOrder, orderBeforeShuffling) {
		t.Error("Unless you are devil unlucky, cards are not shuffled after dealing them all")
	}
}

func TestIfAllCardsAreDealt(t *testing.T) {
	fullDeck := CreateFullDeck()
	cardsMap := map[Card]bool{}
	for range fullDeck.Cards {
		cardsMap[fullDeck.DealNextCard()] = true
	}

	names := [13]CardName{C2, C3, C4, C5, C6, C7, C8, C9, C10, Jack, Queen, King, Ace}
	suits := [4]CardSuit{Hearts, Diamonds, Spades, Clubs}

	for _, suit := range suits {
		for _, name := range names {
			if _, ok := cardsMap[Card{name, suit}]; !ok {
				t.Errorf("%v has not been dealt but is supposed to", Card{name, suit})
			}
		}
	}

}
