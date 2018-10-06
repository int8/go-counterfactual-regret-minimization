package gocfr

import (
	"reflect"
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

	for i := range fullDeck.cards {
		if fullDeck.CardsLeft() != 52-uint8(i) {
			t.Errorf("Full deck should have %v cards left dealing %v card but have %v", 52-i, i, fullDeck.CardsLeft())
		}
		fullDeck.DealNextCard()
	}

}

func TestFullDeckCardsShuffling(t *testing.T) {
	fullDeck := CreateFullDeck(true)
	orderBeforeShuffling := make([]uint8, 52, 52)
	copy(orderBeforeShuffling, fullDeck.shuffleOrder)

	for range fullDeck.cards {
		fullDeck.DealNextCard()
	}

	if reflect.DeepEqual(fullDeck.shuffleOrder, orderBeforeShuffling) {
		t.Error("Unless you are ultra unlucky, cards are not shuffled after dealing them all")
	}
}

func TestIfAllCardsAreDealt(t *testing.T) {
	fullDeck := CreateFullDeck(true)
	cardsMap := map[Card]bool{}
	for range fullDeck.cards {
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

func TestDeckPreparationForTestingPurposes(t *testing.T) {
	aceHearts := Card{Ace, Hearts}
	twoSpades := Card{C2, Spades}
	jackHearts := Card{Jack, Hearts}
	kingHearts := Card{King, Hearts}

	deck := prepareDeckForTest(aceHearts, twoSpades, jackHearts, kingHearts)
	dealtCard := deck.DealNextCard()
	if dealtCard != aceHearts {
		t.Errorf("%v should be dealt but %v was dealt instead", aceHearts, dealtCard)
	}

	dealtCard = deck.DealNextCard()
	if dealtCard != twoSpades {
		t.Errorf("%v should be dealt but %v was dealt instead", twoSpades, dealtCard)
	}

	dealtCard = deck.DealNextCard()
	if dealtCard != jackHearts {
		t.Errorf("%v should be dealt but %v was dealt instead", jackHearts, dealtCard)
	}

	dealtCard = deck.DealNextCard()
	if dealtCard != kingHearts {
		t.Errorf("%v should be dealt but %v was dealt instead", kingHearts, dealtCard)
	}
}
