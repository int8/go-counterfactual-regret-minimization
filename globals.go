package gopoker

const (
	PlayerA  ActorId = 1
	PlayerB          = -PlayerA
	ChanceId         = 0
)

const NoCardSuit CardSuit = 0
const (
	Hearts   CardSuit = 1 + iota // ♥
	Diamonds                     // ♦
	Spades                       // ♠
	Clubs                        // ♣
)

const NoCardName CardName = 0

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
	Start Round = iota
	PreFlop
	Flop
	Turn
	River
	End
)

const (
	DealPublicCards = 1 + iota
	DealPrivateCards
	Fold
	Check
	Bet
	Call
	Raise
)

const InformationSetSize int = 64

var C2Hearts = Card{C2, Hearts}
var C3Hearts = Card{C3, Hearts}
var C4Hearts = Card{C4, Hearts}
var C5Hearts = Card{C5, Hearts}
var C6Hearts = Card{C6, Hearts}
var C7Hearts = Card{C7, Hearts}
var C8Hearts = Card{C8, Hearts}
var C9Hearts = Card{C9, Hearts}
var C10Hearts = Card{C10, Hearts}
var JackHearts = Card{Jack, Hearts}
var QueenHearts = Card{Queen, Hearts}
var KingHearts = Card{King, Hearts}
var AceHearts = Card{Ace, Hearts}

var C2Spades = Card{C2, Spades}
var C3Spades = Card{C3, Spades}
var C4Spades = Card{C4, Spades}
var C5Spades = Card{C5, Spades}
var C6Spades = Card{C6, Spades}
var C7Spades = Card{C7, Spades}
var C8Spades = Card{C8, Spades}
var C9Spades = Card{C9, Spades}
var C10Spades = Card{C10, Spades}
var JackSpades = Card{Jack, Spades}
var QueenSpades = Card{Queen, Spades}
var KingSpades = Card{King, Spades}
var AceSpades = Card{Ace, Spades}

var C2Diamonds = Card{C2, Diamonds}
var C3Diamonds = Card{C3, Diamonds}
var C4Diamonds = Card{C4, Diamonds}
var C5Diamonds = Card{C5, Diamonds}
var C6Diamonds = Card{C6, Diamonds}
var C7Diamonds = Card{C7, Diamonds}
var C8Diamonds = Card{C8, Diamonds}
var C9Diamonds = Card{C9, Diamonds}
var C10Diamonds = Card{C10, Diamonds}
var JackDiamonds = Card{Jack, Diamonds}
var QueenDiamonds = Card{Queen, Diamonds}
var KingDiamonds = Card{King, Diamonds}
var AceDiamonds = Card{Ace, Diamonds}

var C2Clubs = Card{C2, Clubs}
var C3Clubs = Card{C3, Clubs}
var C4Clubs = Card{C4, Clubs}
var C5Clubs = Card{C5, Clubs}
var C6Clubs = Card{C6, Clubs}
var C7Clubs = Card{C7, Clubs}
var C8Clubs = Card{C8, Clubs}
var C9Clubs = Card{C9, Clubs}
var C10Clubs = Card{C10, Clubs}
var JackClubs = Card{Jack, Clubs}
var QueenClubs = Card{Queen, Clubs}
var KingClubs = Card{King, Clubs}
var AceClubs = Card{Ace, Clubs}

var allCards = []*Card{&C2Hearts, &C3Hearts, &C4Hearts, &C5Hearts, &C6Hearts, &C7Hearts, &C8Hearts,
	&C9Hearts, &C10Hearts, &JackHearts, &QueenHearts, &KingHearts, &AceHearts,
	&C2Spades, &C3Spades, &C4Spades, &C5Spades, &C6Spades, &C7Spades, &C8Spades,
	&C9Spades, &C10Spades, &JackSpades, &QueenSpades, &KingSpades, &AceSpades,
	&C2Clubs, &C3Clubs, &C4Clubs, &C5Clubs, &C6Clubs, &C7Clubs, &C8Clubs,
	&C9Clubs, &C10Clubs, &JackClubs, &QueenClubs, &KingClubs, &AceClubs,
	&C2Diamonds, &C3Diamonds, &C4Diamonds, &C5Diamonds, &C6Diamonds, &C7Diamonds, &C8Diamonds,
	&C9Diamonds, &C10Diamonds, &JackDiamonds, &QueenDiamonds, &KingDiamonds, &AceDiamonds,
}
