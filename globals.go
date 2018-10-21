package gopoker

const (
	PlayerA  ActorID = 1
	PlayerB          = -PlayerA
	ChanceId         = 0
)

var (
	NoCardSuit CardSuit = to3BinArray(0) // no card
	Hearts     CardSuit = to3BinArray(1) // ♥
	Diamonds            = to3BinArray(2) // ♦
	Spades              = to3BinArray(3) // ♠
	Clubs               = to3BinArray(4) // ♣
)

var (
	NoCardSymbol CardSymbol = to4BinArray(0) // no card
	C2           CardSymbol = to4BinArray(1)
	C3           CardSymbol = to4BinArray(2)
	C4           CardSymbol = to4BinArray(3)
	C5           CardSymbol = to4BinArray(4)
	C6           CardSymbol = to4BinArray(5)
	C7           CardSymbol = to4BinArray(6)
	C8           CardSymbol = to4BinArray(7)
	C9           CardSymbol = to4BinArray(8)
	C10          CardSymbol = to4BinArray(9)
	Jack         CardSymbol = to4BinArray(10)
	Queen        CardSymbol = to4BinArray(11)
	King         CardSymbol = to4BinArray(12)
	Ace          CardSymbol = to4BinArray(13)
)

const (
	Start Round = iota
	PreFlop
	Flop
	Turn
	River
	End
)

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

var (
	C2Hearts    = Card{C2, Hearts}
	C3Hearts    = Card{C3, Hearts}
	C4Hearts    = Card{C4, Hearts}
	C5Hearts    = Card{C5, Hearts}
	C6Hearts    = Card{C6, Hearts}
	C7Hearts    = Card{C7, Hearts}
	C8Hearts    = Card{C8, Hearts}
	C9Hearts    = Card{C9, Hearts}
	C10Hearts   = Card{C10, Hearts}
	JackHearts  = Card{Jack, Hearts}
	QueenHearts = Card{Queen, Hearts}
	KingHearts  = Card{King, Hearts}
	AceHearts   = Card{Ace, Hearts}

	C2Spades    = Card{C2, Spades}
	C3Spades    = Card{C3, Spades}
	C4Spades    = Card{C4, Spades}
	C5Spades    = Card{C5, Spades}
	C6Spades    = Card{C6, Spades}
	C7Spades    = Card{C7, Spades}
	C8Spades    = Card{C8, Spades}
	C9Spades    = Card{C9, Spades}
	C10Spades   = Card{C10, Spades}
	JackSpades  = Card{Jack, Spades}
	QueenSpades = Card{Queen, Spades}
	KingSpades  = Card{King, Spades}
	AceSpades   = Card{Ace, Spades}

	C2Diamonds    = Card{C2, Diamonds}
	C3Diamonds    = Card{C3, Diamonds}
	C4Diamonds    = Card{C4, Diamonds}
	C5Diamonds    = Card{C5, Diamonds}
	C6Diamonds    = Card{C6, Diamonds}
	C7Diamonds    = Card{C7, Diamonds}
	C8Diamonds    = Card{C8, Diamonds}
	C9Diamonds    = Card{C9, Diamonds}
	C10Diamonds   = Card{C10, Diamonds}
	JackDiamonds  = Card{Jack, Diamonds}
	QueenDiamonds = Card{Queen, Diamonds}
	KingDiamonds  = Card{King, Diamonds}
	AceDiamonds   = Card{Ace, Diamonds}

	C2Clubs    = Card{C2, Clubs}
	C3Clubs    = Card{C3, Clubs}
	C4Clubs    = Card{C4, Clubs}
	C5Clubs    = Card{C5, Clubs}
	C6Clubs    = Card{C6, Clubs}
	C7Clubs    = Card{C7, Clubs}
	C8Clubs    = Card{C8, Clubs}
	C9Clubs    = Card{C9, Clubs}
	C10Clubs   = Card{C10, Clubs}
	JackClubs  = Card{Jack, Clubs}
	QueenClubs = Card{Queen, Clubs}
	KingClubs  = Card{King, Clubs}
	AceClubs   = Card{Ace, Clubs}
	NoCard     = Card{NoCardSymbol, NoCardSuit}
)

var allCards = []*Card{&C2Hearts, &C3Hearts, &C4Hearts, &C5Hearts, &C6Hearts, &C7Hearts, &C8Hearts,
	&C9Hearts, &C10Hearts, &JackHearts, &QueenHearts, &KingHearts, &AceHearts,
	&C2Spades, &C3Spades, &C4Spades, &C5Spades, &C6Spades, &C7Spades, &C8Spades,
	&C9Spades, &C10Spades, &JackSpades, &QueenSpades, &KingSpades, &AceSpades,
	&C2Clubs, &C3Clubs, &C4Clubs, &C5Clubs, &C6Clubs, &C7Clubs, &C8Clubs,
	&C9Clubs, &C10Clubs, &JackClubs, &QueenClubs, &KingClubs, &AceClubs,
	&C2Diamonds, &C3Diamonds, &C4Diamonds, &C5Diamonds, &C6Diamonds, &C7Diamonds, &C8Diamonds,
	&C9Diamonds, &C10Diamonds, &JackDiamonds, &QueenDiamonds, &KingDiamonds, &AceDiamonds,
}
