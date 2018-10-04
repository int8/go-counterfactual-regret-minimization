package gocfr

type PokerTable interface {
	AddToPot(amount float64)
	DropPublicCard(card Card)
}

type RhodeIslandPokerTable struct {
	potSize float64
	publicCards []Card
}

func (table *RhodeIslandPokerTable) AddToPot(amount float64) {
	table.potSize += amount
}

func (table *RhodeIslandPokerTable) DropPublicCard(card Card) {
	table.publicCards = append(table.publicCards, card)
}

