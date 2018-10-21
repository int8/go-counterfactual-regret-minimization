package gopoker

type ActionName [3]bool

type Action interface {
	Name() ActionName
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
	case DealPublicCards:
		return "DealPublicCards"
	}
	return "Undefined"
}
