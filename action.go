package gocfr

type Action byte

const (
	NoAction Action = iota
	DealPublicCard
	DealPrivateCards
	Fold
	Check
	Bet
	Call
	Raise
)

const MaxRaises int = 3

const PreFlopBetSize float64 = 10.
const PostFlopBetSize float64 = 20.
const Ante float64 = 5.0

func (m Action) String() string {
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
