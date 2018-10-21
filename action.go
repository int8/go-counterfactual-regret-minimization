package gopoker

type ActionName [3]bool

type Action interface {
	Name() ActionName
}

func (m ActionName) String() string {
	switch m {
	case Check:
		return "Ch"
	case Bet:
		return "B"
	case Call:
		return "C"
	case Fold:
		return "F"
	case Raise:
		return "R"
	case DealPrivateCards:
		return "DPrv"
	case DealPublicCards:
		return "DPub"
	}
	return "?"
}
