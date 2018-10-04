package gocfr

type Player int8
const (
	Environment  Player = 0
	PlayerA Player = 1
	PlayerB        = -PlayerA
)

type HeadsUpPokerPlayer interface {
	Opponent() HeadsUpPokerPlayer
	CollectPrivateCard(card Card)
}

type RhodeIslandPokerPlayer struct {
	privateCard Card
	player Player
	stack float64
	opponent *RhodeIslandPokerPlayer
}

func (player *RhodeIslandPokerPlayer) Opponent() *RhodeIslandPokerPlayer{
	return player.opponent
}

func (player *RhodeIslandPokerPlayer) CollectPrivateCard(card Card) {
	player.privateCard = card
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



