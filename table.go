package gocfr

type Table struct {
	pot   float64
	cards []Card
}

func (table *Table) AddToPot(amount float64) {
	table.pot += amount
}

func (table *Table) DropPublicCard(card *Card) {
	table.cards = append(table.cards, *card)
}

func (table *Table) Clone() *Table {
	cards := make([]Card, len(table.cards))
	copy(cards, table.cards)
	return &Table{pot: table.pot, cards: cards}
}
