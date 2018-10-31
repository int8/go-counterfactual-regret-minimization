package rounds

type PokerRound int8

const (
	Start PokerRound = iota
	PreFlop
	Flop
	Turn
	River
	End
)

func (round PokerRound) NextRound() PokerRound {
	switch round {
	case Start:
		return PreFlop
	case PreFlop:
		return Flop
	case Flop:
		return Turn
	case Turn:
		return River
	}

	return End
}

func (round PokerRound) String() string {
	switch round {
	case Start:
		return "Start"
	case End:
		return "End"
	case PreFlop:
		return "Preflop"
	case Turn:
		return "Turn"
	case River:
		return "River"
	case Flop:
		return "Flop"
	}
	return "(?)"
}
