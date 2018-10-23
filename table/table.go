package table

import "github.com/int8/gopoker/cards"

type PokerTable struct {
	Pot   float32
	Cards []cards.Card
}

func (table *PokerTable) AddToPot(amount float32) {
	table.Pot += amount
}

func (table *PokerTable) DropPublicCard(card *cards.Card) {
	table.Cards = append(table.Cards, *card)
}

func (table *PokerTable) Clone() *PokerTable {
	tablecards := make([]cards.Card, len(table.Cards))
	copy(tablecards, table.Cards)
	return &PokerTable{Pot: table.Pot, Cards: tablecards}
}
