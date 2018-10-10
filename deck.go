package gocfr

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type CardName byte

type CardSuit byte

type Card struct {
	name CardName
	suit CardSuit
}

type FullDeck struct {
	cards map[*Card]bool
}

func CreateFullDeck(shuffleInitially bool) *FullDeck {

	fullDeck := *new(FullDeck)
	fullDeck.cards = make(map[*Card]bool, 52)
	for _, card := range allCards {
		fullDeck.cards[card] = true
	}
	fullDeck.Shuffle()
	return &fullDeck
}

func (d *FullDeck) Shuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (d *FullDeck) DealNextRandomCard() *Card {
	var card *Card
	i := rand.Intn(len(d.cards))
	for card = range d.cards {
		if i == 0 {
			break
		}
		i--
	}
	delete(d.cards, card)
	return card
}

func (d *FullDeck) CardsLeft() int {
	return len(d.cards)
}

func (d *FullDeck) RemainingCards() []*Card {
	cards := make([]*Card, 0, len(d.cards))
	for card := range d.cards {
		cards = append(cards, card)
	}
	return cards
}

func (d *FullDeck) Clone() *FullDeck {
	cards := make(map[*Card]bool, len(d.cards))
	for k := range d.cards {
		cards[k] = true
	}
	return &FullDeck{cards}
}

func (c Card) String() string {
	return fmt.Sprintf("%v%v", c.suit, c.name)
}

func (s CardSuit) String() string {
	switch s {
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Spades:
		return "♠"
	case Clubs:
		return "♣"
	}
	return "?"
}

func (n CardName) String() string {
	switch n {
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		return strconv.Itoa(int(n))
	}
}
