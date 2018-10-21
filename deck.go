package gopoker

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Deck interface {
	Shuffle()
	RemoveCard(card *Card)
	CardsLeft() int
	Clone() Deck
	RemainingCards() []*Card
}

// CardName - there is 1-10 + J Q K A = 14 cards + NoCardName identifier / 4 bits is enough
type CardName [4]bool

// CardSuit - there is 4 card suits + NoCardSuit identifier - 3 bits is ok
type CardSuit [3]bool

type Card struct {
	Name CardName
	Suit CardSuit
}

type FullDeck struct {
	Cards map[*Card]bool
}

func CreateFullDeck(shuffleInitially bool) *FullDeck {

	fullDeck := *new(FullDeck)
	fullDeck.Cards = make(map[*Card]bool, 52)
	for _, card := range allCards {
		fullDeck.Cards[card] = true
	}
	fullDeck.Shuffle()

	return &fullDeck
}

func (d *FullDeck) Shuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (d *FullDeck) RemoveCard(card *Card) {
	delete(d.Cards, card)
}

func (d *FullDeck) CardsLeft() int {
	return len(d.Cards)
}

func (d *FullDeck) RemainingCards() []*Card {
	cards := make([]*Card, 0, len(d.Cards))
	for card := range d.Cards {
		cards = append(cards, card)
	}
	return cards
}

func (d *FullDeck) Clone() Deck {
	cards := make(map[*Card]bool, len(d.Cards))
	for k := range d.Cards {
		cards[k] = true
	}
	return &FullDeck{cards}
}

// TODO: what happens when no cards are left ?
func (d *FullDeck) DealNextRandomCard() *Card {
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

type LimitedDeck struct {
	Cards map[*Card]bool
}

func CreateLimitedDeck(minCardName CardName, shuffleInitially bool) *LimitedDeck {

	deck := *new(LimitedDeck)
	deck.Cards = make(map[*Card]bool, 20)
	for _, card := range allCards {
		if cardNameCompare(card.Name, minCardName) > 0 {
			deck.Cards[card] = true
		}
	}
	deck.Shuffle()

	return &deck
}

func (d *LimitedDeck) Shuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (d *LimitedDeck) RemoveCard(card *Card) {
	delete(d.Cards, card)
}

func (d *LimitedDeck) CardsLeft() int {
	return len(d.Cards)
}

func (d *LimitedDeck) RemainingCards() []*Card {
	cards := make([]*Card, 0, len(d.Cards))
	for card := range d.Cards {
		cards = append(cards, card)
	}
	return cards
}

func (d *LimitedDeck) Clone() Deck {
	cards := make(map[*Card]bool, len(d.Cards))
	for k := range d.Cards {
		cards[k] = true
	}
	return &LimitedDeck{cards}
}

func (c Card) String() string {
	return fmt.Sprintf("%v%v", c.Suit, c.Name)
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
		return strconv.Itoa(int(CardNameInt(n)))
	}
}
