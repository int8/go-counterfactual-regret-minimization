package kuhn

import (
	. "github.com/int8/gopoker"
	"math/rand"
	"time"
)

type KuhnDeck struct {
	Cards map[*Card]bool
}

func CreateKuhnDeck(shuffleInitially bool) *KuhnDeck {

	deck := *new(KuhnDeck)
	deck.Cards = make(map[*Card]bool, 3)
	deck.Cards[&JackHearts] = true
	deck.Cards[&QueenHearts] = true
	deck.Cards[&KingHearts] = true
	deck.Shuffle()

	return &deck
}

func (d *KuhnDeck) Shuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (d *KuhnDeck) RemoveCard(card *Card) {
	delete(d.Cards, card)
}

func (d *KuhnDeck) CardsLeft() int {
	return len(d.Cards)
}

func (d *KuhnDeck) RemainingCards() []*Card {
	cards := make([]*Card, 0, len(d.Cards))
	for card := range d.Cards {
		cards = append(cards, card)
	}
	return cards
}

func (d *KuhnDeck) Clone() Deck {
	cards := make(map[*Card]bool, len(d.Cards))
	for k := range d.Cards {
		cards[k] = true
	}
	return &FullDeck{cards}
}

func (d *KuhnDeck) DealNextRandomCard() *Card {
	var card *Card
	i := rand.Intn(len(d.Cards))
	for card = range d.Cards {
		if i == 0 {
			break
		}
		i--
	}
	d.RemoveCard(card)
	return card
}
