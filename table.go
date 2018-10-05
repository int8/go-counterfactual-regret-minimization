package gocfr


type PokerTable struct {
	potSize float64
	publicCards []Card
}

func (table *PokerTable) AddToPot(amount float64) {
	table.potSize += amount
}

func (table *PokerTable) DropPublicCard(card Card) {
	table.publicCards = append(table.publicCards, card)
}

