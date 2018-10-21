package gopoker

import (
	"testing"
)

// TODO: DealNextRandomCard is not used
func TestFullDeckCardsCount(t *testing.T) {
	fullDeck := CreateFullDeck(true)

	if len(fullDeck.Cards) != 52 {
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

func TestLimitedDeckCardsCount(t *testing.T) {
	deck := CreateLimitedDeck(C10, true)

	if len(deck.Cards) != 20 {
		t.Error("Limited deck starting from 10 should count 20 cards")
	}

	if deck.CardsLeft() != 20 {
		t.Errorf("Limited deck starting from 10 should have 20 cards left after initialization but have %v", deck.CardsLeft())
	}
}

// TODO: DealNextRandomCard is not used
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
