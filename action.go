package gocfr

type ActionName byte

type Action interface {
	Name() ActionName
}

type PlayerAction struct {
	name ActionName
}

func (a PlayerAction) Name() ActionName {
	return a.name
}

type DealPrivateCardsAction struct {
	cardA *Card
	cardB *Card
}

func (a DealPrivateCardsAction) Name() ActionName {
	return DealPrivateCards
}

type DealPublicCardAction struct {
	card *Card
}

func (a DealPublicCardAction) Name() ActionName {
	return DealPublicCard
}

func (m ActionName) String() string {
	switch m {
	case Check:
		return "Check"
	case Bet:
		return "Bet"
	case Call:
		return "Call"
	case Fold:
		return "Fold"
	case Raise:
		return "Raise"
	case DealPrivateCards:
		return "DealPrivateCards"
	case DealPublicCard:
		return "DealPublicCard"
	}
	return "Undefined"
}
