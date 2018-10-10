package gocfr

type Round int8

func (round Round) NextRound() Round {
	switch round {
	case Start:
		return PreFlop
	case PreFlop:
		return Flop
	case Flop:
		return Turn
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
	case Flop:
		return "Flop"
	}
	return "(?)"
}
