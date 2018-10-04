package gocfr

type Player int8
const (
	Environment  Player = 0
	PlayerA Player = 1
	PlayerB        = -PlayerA
)

type HeadsUpPokerPlayer struct {
	player Player
	stack float64
	opponent *HeadsUpPokerPlayer
}

func (p Player) String() string {
	if p == PlayerA {
		return "A"
	}
	if p == PlayerB {
		return "B"
	}
	return "Chance"
}



