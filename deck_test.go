package gocfr

import (
	"testing"
)

func TestFullDeckCardsCount(t *testing.T) {
	fullDeck := CreateFullDeck(true)

	if len(fullDeck.cards) != 52 {
		t.Error("Full deck should count 52 cards")
	}

	if fullDeck.CardsLeft() != 52 {
		t.Errorf("Full deck should have 52 cards left after initialization but have %v", fullDeck.CardsLeft())
	}

	for i := range fullDeck.RemainingCards() {
		if fullDeck.CardsLeft() != 52-i {
			t.Errorf("Full deck should have %v cards left dealing %v card but have %v", 52-i, i, fullDeck.CardsLeft())
		}
		fullDeck.DealNextRandomCard()
	}

}

func TestIfAllCardsAreDealt(t *testing.T) {
	fullDeck := CreateFullDeck(true)
	cardsMap := map[Card]bool{}
	deckSize := fullDeck.CardsLeft()
	for i := 0; i < deckSize; i++ {
		cardsMap[*fullDeck.DealNextRandomCard()] = true
	}

	for _, card := range allCards {
		if _, ok := cardsMap[*card]; !ok {
			t.Errorf("%v has not been dealt but is supposed to", card)
		}
	}
}
