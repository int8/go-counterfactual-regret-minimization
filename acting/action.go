package acting

type ActionName [3]bool

type Action interface {
	Name() ActionName
}

var (
	NoAction         ActionName = to3BinArray(0)
	DealPublicCards  ActionName = to3BinArray(1)
	DealPrivateCards ActionName = to3BinArray(2)
	Fold             ActionName = to3BinArray(3)
	Check            ActionName = to3BinArray(4)
	Bet              ActionName = to3BinArray(5)
	Call             ActionName = to3BinArray(6)
	Raise            ActionName = to3BinArray(7)
)

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
