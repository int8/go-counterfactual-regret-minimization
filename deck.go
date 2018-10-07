package gocfr

import (
	"fmt"
	"math/rand"
	"strconv"
)

type CardName uint8

const (
	Hearts   CardSuit = iota // ♥
	Diamonds                 // ♦
	Spades                   // ♠
	Clubs                    // ♣
)

type CardSuit uint8

const (
	C2 CardName = 2 + iota
	C3
	C4
	C5
	C6
	C7
	C8
	C9
	C10
	Jack
	Queen
	King
	Ace
)

type Card struct {
	name CardName
	suit CardSuit
}

type FullDeck struct {
	cards            []Card
	shuffleOrder     []uint8
	currentCardIndex uint8
}

func CreateFullDeck(shuffleInitially bool) *FullDeck {
	names := [13]CardName{C2, C3, C4, C5, C6, C7, C8, C9, C10, Jack, Queen, King, Ace}
	suits := [4]CardSuit{Hearts, Diamonds, Spades, Clubs}
	fullDeck := *new(FullDeck)

	for _, suit := range suits {
		for _, name := range names {
			fullDeck.cards = append(fullDeck.cards, Card{name, suit})
		}
	}
	fullDeck.shuffleOrder = makeRange(0, 51)
	if shuffleInitially {
		fullDeck.Shuffle()
	}
	fullDeck.currentCardIndex = 0
	return &fullDeck
}

func (d *FullDeck) Shuffle() {
	offset := d.currentCardIndex
	order := d.shuffleOrder
	rand.Shuffle(51-int(offset), func(i int, j int) {
		order[int(offset)+i], order[int(offset)+j] = order[int(offset)+j], order[int(offset)+i]
	})
}

func (d *FullDeck) DealNextCard() *Card {

	cardToBeReturned := d.cards[d.shuffleOrder[d.currentCardIndex]]
	d.currentCardIndex = (d.currentCardIndex + 1) % 52
	// if all cards dealt - shuffle
	if d.currentCardIndex == 0 {
		d.Shuffle()
	}
	return &cardToBeReturned
}

func (d *FullDeck) CardsLeft() uint8 {
	return 52 - d.currentCardIndex
}

func (d *FullDeck) Clone() *FullDeck {
	shuffleOrder := make([]uint8, len(d.shuffleOrder))
	copy(shuffleOrder, d.shuffleOrder)
	// important to reuse d.cards here
	return &FullDeck{d.cards, shuffleOrder, d.currentCardIndex}
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
