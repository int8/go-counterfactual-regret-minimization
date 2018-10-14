package gopoker

type Table struct {
	Pot   float32
	Cards []Card
}

func (table *Table) AddToPot(amount float32) {
	table.Pot += amount
}

func (table *Table) DropPublicCard(card *Card) {
	table.Cards = append(table.Cards, *card)
}

func (table *Table) Clone() *Table {
	cards := make([]Card, len(table.Cards))
	copy(cards, table.Cards)
	return &Table{Pot: table.Pot, Cards: cards}
}
