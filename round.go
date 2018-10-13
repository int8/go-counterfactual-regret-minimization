package gopoker

type Round int8

func (round Round) NextRound() Round {
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

func (round Round) String() string {
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
