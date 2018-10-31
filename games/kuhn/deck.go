package kuhn

import (
	"github.com/int8/gopoker/cards"
	"math/rand"
	"time"
)

type KuhnDeck struct {
	Cards map[*cards.Card]bool
}

func CreateKuhnDeck() *KuhnDeck {

	deck := *new(KuhnDeck)
	deck.Cards = make(map[*cards.Card]bool, 3)
	deck.Cards[&cards.JackHearts] = true
	deck.Cards[&cards.QueenHearts] = true
	deck.Cards[&cards.KingHearts] = true
	deck.Shuffle()

	return &deck
}

func (d *KuhnDeck) Shuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (d *KuhnDeck) RemoveCard(card *cards.Card) {
	delete(d.Cards, card)
}

func (d *KuhnDeck) CardsLeft() int {
	return len(d.Cards)
}

func (d *KuhnDeck) RemainingCards() []*cards.Card {
	kuhncards := make([]*cards.Card, 0, len(d.Cards))
	for card := range d.Cards {
		kuhncards = append(kuhncards, card)
	}
	return kuhncards
}

func (d *KuhnDeck) Clone() cards.Deck {
	kuhncards := make(map[*cards.Card]bool, len(d.Cards))
	for k := range d.Cards {
		kuhncards[k] = true
	}
	return &KuhnDeck{kuhncards}
}

func (d *KuhnDeck) DealNextRandomCard() *cards.Card {
	var card *cards.Card
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
