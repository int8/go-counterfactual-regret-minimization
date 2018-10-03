package gocfr

const (
	Hearts CardSuit = iota // ♥
	Diamonds // ♦
	Spades   // ♠
	Clubs    // ♣
)

const (
	C2 CardName = 2 + iota
	C3
	C4
	C5
	C6
	C7
	C8
	C9
	C10
	Jack
	Queen
	King
	Ace
)


const (
	Check Move = iota
	Bet
	Raise
	Call
	Fold
	DealPublicCard
	DealPrivateCards
)

const (
	Chance Player = 0
	PlayerA Player = 1
	PlayerB = - PlayerA

)

const (
	Start Round = iota
	PreFlop
	Flop
	Turn
	End
)

