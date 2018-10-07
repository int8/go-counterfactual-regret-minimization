package gocfr

import (
	"math"
	"testing"
)

type MoveTestsTriple struct {
	move     Move
	preTest  func(state *GameState) bool
	postTest func(state *GameState) bool
}

func TestGameCreation(t *testing.T) {
	root := createRootForTest(100., 100.)
	if root.causingMove != NoMove {
		t.Error("Root node should not have causing action")
	}

	if root.parent != nil {
		t.Error("Root node should not have nil parent")
	}

	if root.round != Start {
		t.Error("Initial round of the game should be Start")
	}

	if root.IsTerminal() == true {
		t.Error("Game root should not be terminal")
	}

	moves := root.actors[root.nextToMove].GetAvailableMoves(root)

	if moves == nil {
		t.Error("Game root should have one action available, no actions available")
	}

	if len(moves) != 1 {
		t.Errorf("Game root should have one action available, %v actions available", len(moves))
	}

}

func TestIfParentsCorrect(t *testing.T) {
	root := createRootForTest(100., 100.)
	child := root.actors[root.nextToMove].(*Chance).Act(root, DealPrivateCards)
	if child.parent != root {
		t.Error("Root child should have root as a parent")
	}
}

func TestIfStackLimitsAvailableActions(t *testing.T) {
	root5 := createRootForTest(5., 5.)
	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, noTest(), noTest()},
		{Check, onlyCheckAvailable(), noTest()},
		{Check, onlyCheckAvailable(), noTest()},
		{DealPublicCard, noTest(), noTest()},
		{Check, onlyCheckAvailable(), noTest()},
		{Check, onlyCheckAvailable(), noTest()},
		{DealPublicCard, noTest(), noTest()},
		{Check, onlyCheckAvailable(), noTest()},
		{Check, onlyCheckAvailable(), noTest()},
	}
	testGamePlayAfterEveryMove(root5, movesTestsPairs, t)

	root15 := createRootForTest(15., 15.)
	movesTestsPairs = []MoveTestsTriple{
		{DealPrivateCards, noTest(), noTest()},
		{Check, checkAndBetAvailable(), noTest()},
		{Check, checkAndBetAvailable(), noTest()},
		{DealPublicCard, noTest(), noTest()},
		{Check, onlyCheckAvailable(), noTest()}, // at this point bet size exceeds players stack
		{Check, onlyCheckAvailable(), noTest()}, // only check available
		{DealPublicCard, noTest(), noTest()},
		{Check, onlyCheckAvailable(), noTest()}, // at this point bet size exceeds players stack
		{Check, onlyCheckAvailable(), noTest()}, // only check available
	}

	testGamePlayAfterEveryMove(root15, movesTestsPairs, t)

	root1000_15 := createRootForTest(1000., 15.)
	movesTestsPairs = []MoveTestsTriple{
		{DealPrivateCards, noTest(), noTest()},
		{Check, checkAndBetAvailable(), noTest()},
		{Check, checkAndBetAvailable(), noTest()},
		{DealPublicCard, noTest(), noTest()},
		{Check, onlyCheckAvailable(), noTest()}, // at this point bet size exceeds one of the players stack
		{Check, onlyCheckAvailable(), noTest()}, // only check available
		{DealPublicCard, noTest(), noTest()},
		{Check, onlyCheckAvailable(), noTest()}, // at this point bet size exceeds one of the players stack
		{Check, onlyCheckAvailable(), noTest()}, // only check available
	}

	testGamePlayAfterEveryMove(root1000_15, movesTestsPairs, t)

	root1000_35 := createRootForTest(1000., 35.)
	movesTestsPairs = []MoveTestsTriple{
		{DealPrivateCards, noTest(), noTest()},
		{Check, checkAndBetAvailable(), noTest()},
		{Check, checkAndBetAvailable(), noTest()},
		{DealPublicCard, noTest(), noTest()},
		{Check, checkAndBetAvailable(), noTest()},
		{Check, checkAndBetAvailable(), noTest()},
		{DealPublicCard, noTest(), noTest()},
		{Check, checkAndBetAvailable(), noTest()},
		{Check, checkAndBetAvailable(), noTest()},
	}

	testGamePlayAfterEveryMove(root1000_35, movesTestsPairs, t)
}

func TestGamePlay_1(t *testing.T) {
	root := createRootForTest(100., 100.)
	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, roundCheck(Start), roundCheck(PreFlop)},
		{Check, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Check, roundCheck(PreFlop), roundCheck(PreFlop)},
		{DealPublicCard, roundCheck(PreFlop), roundCheck(Flop)},
		{Bet, roundCheck(Flop), roundCheck(Flop)},
		{Call, roundCheck(Flop), roundCheck(Flop)},
		{DealPublicCard, roundCheck(Flop), roundCheck(Turn)},
		{Check, roundCheck(Turn), roundCheck(Turn)},
		{Bet, roundCheck(Turn), roundCheck(Turn)},
		{Call, roundCheck(Turn), gameEnd()},
	}

	testGamePlayAfterEveryMove(root, movesTestsPairs, t)
}

func TestGamePlay_MaxRaises(t *testing.T) {
	root := createRootForTest(100., 100.)

	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, roundCheck(Start), roundCheck(PreFlop)},
		{Bet, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), noRaiseAvailable()},
		{Fold, roundCheck(PreFlop), gameEnd()},
	}

	testGamePlayAfterEveryMove(root, movesTestsPairs, t)

	movesTestsPairs = []MoveTestsTriple{
		{DealPrivateCards, roundCheck(Start), roundCheck(PreFlop)},
		{Bet, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), noRaiseAvailable()},
		{Fold, roundCheck(PreFlop), gameEnd()},
	}

	testGamePlayAfterEveryMove(root, movesTestsPairs, t)

	movesTestsPairs = []MoveTestsTriple{
		{DealPrivateCards, roundCheck(Start), roundCheck(PreFlop)},
		{Bet, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), noRaiseAvailable()},
		{Call, roundCheck(PreFlop), roundCheck(PreFlop)},
		{DealPublicCard, roundCheck(PreFlop), roundCheck(Flop)},
		{Bet, roundCheck(Flop), roundCheck(Flop)},
		{Raise, roundCheck(Flop), roundCheck(Flop)},
		{Raise, roundCheck(Flop), roundCheck(Flop)},
		{Call, roundCheck(Flop), roundCheck(Flop)},
		{DealPublicCard, roundCheck(Flop), roundCheck(Turn)},
		{Check, roundCheck(Turn), roundCheck(Turn)},
		{Check, roundCheck(Turn), gameEnd()},
	}

	testGamePlayAfterEveryMove(root, movesTestsPairs, t)

}

func TestGamePlay_CheckIfPlayerToMoveCorrect(t *testing.T) {
	root := createRootForTest(100., 100.)
	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, actorToMove(ChanceId), actorToMove(PlayerA)},
		{Check, actorToMove(PlayerA), actorToMove(PlayerB)},
		{Check, actorToMove(PlayerB), actorToMove(ChanceId)},
		{DealPublicCard, actorToMove(ChanceId), actorToMove(PlayerA)},
		{Bet, actorToMove(PlayerA), actorToMove(PlayerB)},
		{Call, actorToMove(PlayerB), actorToMove(ChanceId)},
		{DealPublicCard, actorToMove(ChanceId), actorToMove(PlayerA)},
		{Check, actorToMove(PlayerA), actorToMove(PlayerB)},
		{Bet, actorToMove(PlayerB), actorToMove(PlayerA)},
		{Call, actorToMove(PlayerA), actorToMove(NoActor)},
	}

	testGamePlayAfterEveryMove(root, movesTestsPairs, t)
}

func TestGamePlay_CheckIfStacksChange(t *testing.T) {
	root := createRootForTest(100., 100.)
	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, stackEqualTo(PlayerA, 100.), stackEqualTo(PlayerA, 100.-Ante)},
		{Check, stackEqualTo(PlayerA, 100.-Ante), stackEqualTo(PlayerA, 100.-Ante)},
		{Check, stackEqualTo(PlayerB, 100.-Ante), stackEqualTo(PlayerB, 100.-Ante)},
		{DealPublicCard, stackEqualTo(PlayerA, 100.-Ante), stackEqualTo(PlayerA, 100.-Ante)},
		{Bet, stackEqualTo(PlayerA, 100.-Ante), stackEqualTo(PlayerA, 100.-Ante-PostFlopBetSize)},
		{Call, stackEqualTo(PlayerB, 100.-Ante), stackEqualTo(PlayerB, 100.-Ante-PostFlopBetSize)},
		{DealPublicCard, stackEqualTo(PlayerB, 100.-Ante-PostFlopBetSize), stackEqualTo(PlayerB, 100.-Ante-PostFlopBetSize)},
		{Check, stackEqualTo(PlayerA, 100.-Ante-PostFlopBetSize), stackEqualTo(PlayerA, 100.-Ante-PostFlopBetSize)},
		{Bet, stackEqualTo(PlayerB, 100.-Ante-PostFlopBetSize), stackEqualTo(PlayerB, 100.-Ante-2*PostFlopBetSize)},
		{Call, stackEqualTo(PlayerA, 100.-Ante-PostFlopBetSize), stackEqualTo(PlayerA, 100.-Ante-2*PostFlopBetSize)},
	}

	testGamePlayAfterEveryMove(root, movesTestsPairs, t)
}

func TestGamePlay_CheckIfPotChanges(t *testing.T) {
	root := createRootForTest(100., 100.)
	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, potEqualsTo(0.0), potEqualsTo(10.0)},
		{Check, potEqualsTo(10), potEqualsTo(10)},
		{Check, potEqualsTo(10), potEqualsTo(10)},
		{DealPublicCard, potEqualsTo(10), potEqualsTo(10)},
		{Bet, potEqualsTo(10), potEqualsTo(10 + PostFlopBetSize)},
		{Call, potEqualsTo(10 + PostFlopBetSize), potEqualsTo(10 + 2*PostFlopBetSize)},
		{DealPublicCard, potEqualsTo(10 + 2*PostFlopBetSize), potEqualsTo(10 + 2*PostFlopBetSize)},
		{Check, potEqualsTo(10 + 2*PostFlopBetSize), potEqualsTo(10 + 2*PostFlopBetSize)},
		{Bet, potEqualsTo(10 + 2*PostFlopBetSize), potEqualsTo(10 + 3*PostFlopBetSize)},
		{Call, potEqualsTo(10 + 3*PostFlopBetSize), potEqualsTo(10 + 4*PostFlopBetSize)},
	}

	testGamePlayAfterEveryMove(root, movesTestsPairs, t)
}

func TestIfRootCreationWithDeckPreparedWorks(t *testing.T) {
	aceHearts := Card{Ace, Hearts}
	c2Spades := Card{C2, Spades}
	jackHearts := Card{Jack, Hearts}
	kingHearts := Card{King, Hearts}

	preparedDeck := prepareDeckForTest(aceHearts, c2Spades, jackHearts, kingHearts)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, noTest(), privateCards(aceHearts, c2Spades)},
		{Check, noTest(), noTest()},
		{Check, noTest(), noTest()},
		{DealPublicCard, noTest(), flopCard(jackHearts)},
		{Check, noTest(), noTest()},
		{Check, noTest(), noTest()},
		{DealPublicCard, noTest(), turnCard(kingHearts)},
		{Check, noTest(), noTest()},
		{Check, noTest(), noTest()},
	}

	testGamePlayAfterEveryMove(root, movesTestsPairs, t)

}

func TestGamePlay_CheckIfChildPointersDifferFromParentsPointers(t *testing.T) {
	root := createRootForTest(100., 100.)
	child := root.actors[root.nextToMove].(*Chance).Act(root, DealPrivateCards)

	if child.actors[ChanceId] == root.actors[ChanceId] {
		t.Error("chance actor refers to the same value in both child and parent")
	}

	if child.actors[PlayerA] == root.actors[PlayerA] {
		t.Error("PlayerA actor refers to the same value in both child and parent")
	}

	if child.actors[PlayerB] == root.actors[PlayerB] {
		t.Error("PlayerB actor refers to the same value in both child and parent")
	}

	if child.table == root.table {
		t.Error("table should be different for child and parent")
	}
}

func TestGamePlayEvaluationFlushVsNothing(t *testing.T) {
	privateACard := Card{Ace, Hearts}
	privateBCard := Card{C2, Spades}
	flopPublicCard := Card{Jack, Hearts}
	turnPublicCard := Card{King, Hearts}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Check, Check, DealPublicCard, Bet, Call, DealPublicCard, Check, Bet, Call}
	singlePlayerPotContribution := Ante + PostFlopBetSize + PostFlopBetSize
	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationFlushVsStraightFlush(t *testing.T) {
	privateACard := Card{Ace, Hearts}
	privateBCard := Card{Queen, Hearts}
	flopPublicCard := Card{Jack, Hearts}
	turnPublicCard := Card{King, Hearts}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Bet, Call, DealPublicCard, Bet, Call, DealPublicCard, Check, Bet, Call}
	singlePlayerPotContribution := Ante + PreFlopBetSize + PostFlopBetSize + PostFlopBetSize
	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationPairVsPairDraw(t *testing.T) {
	privateACard := Card{Ace, Hearts}
	privateBCard := Card{Ace, Spades}
	flopPublicCard := Card{Jack, Spades}
	turnPublicCard := Card{King, Hearts}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Check, Check, DealPublicCard, Check, Check, DealPublicCard, Check, Check}
	singlePlayerPotContribution := Ante

	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(0), t)
}

func TestGamePlayEvaluationPairVsPairAWins(t *testing.T) {
	privateACard := Card{Ace, Hearts}
	privateBCard := Card{King, Spades}
	flopPublicCard := Card{Ace, Spades}
	turnPublicCard := Card{King, Hearts}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Check, Check, DealPublicCard, Check, Check, DealPublicCard, Check, Check}
	singlePlayerPotContribution := Ante
	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationPairVsPairBWinsBetterOwnCard(t *testing.T) {
	privateACard := Card{Jack, Hearts}
	privateBCard := Card{King, Spades}
	flopPublicCard := Card{C2, Spades}
	turnPublicCard := Card{C2, Hearts}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Check, Check, DealPublicCard, Check, Check, DealPublicCard, Check, Check}
	singlePlayerPotContribution := Ante
	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationStraightVsStraightAWinsBetterOwnCard(t *testing.T) {
	privateACard := Card{King, Hearts}
	privateBCard := Card{C10, Spades}
	flopPublicCard := Card{Jack, Clubs}
	turnPublicCard := Card{Queen, Diamonds}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Bet, Call, DealPublicCard, Check, Check, DealPublicCard, Check, Check}
	singlePlayerPotContribution := Ante + PreFlopBetSize
	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationPairVsThreeOfAKindBWins(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Bet, Call, DealPublicCard, Check, Check, DealPublicCard, Check, Check}
	singlePlayerPotContribution := Ante + PreFlopBetSize
	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationOwnCardVsOwnCardAWins(t *testing.T) {
	privateACard := Card{King, Hearts}
	privateBCard := Card{C10, Clubs}
	flopPublicCard := Card{C2, Diamonds}
	turnPublicCard := Card{C7, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Bet, Call, DealPublicCard, Check, Check, DealPublicCard, Bet, Call}
	singlePlayerPotContribution := Ante + PreFlopBetSize + PostFlopBetSize
	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAFoldsOnTurn(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Bet, Call, DealPublicCard, Check, Check, DealPublicCard, Check, Bet, Fold}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsOnTurn(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Bet, Call, DealPublicCard, Check, Check, DealPublicCard, Bet, Fold}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAFoldsOnFlop(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Bet, Call, DealPublicCard, Check, Bet, Fold}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsOnFlop(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Bet, Call, DealPublicCard, Bet, Fold}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAFoldsPreFlop(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Check, Bet, Fold}
	singlePlayerPotContribution := Ante

	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsPreFlop(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Check, Bet, Raise, Fold}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsPreFlopManyRaises(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	moves := []Move{DealPrivateCards, Check, Bet, Raise, Raise, Raise, Fold}
	singlePlayerPotContribution := Ante + 3*PreFlopBetSize

	testGamePlayAfterAllMoves(root, moves, gameEnd(), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllMoves(root, moves, gameResult(2*singlePlayerPotContribution), t)
}

func testGamePlayAfterEveryMove(node *GameState, movesTests []MoveTestsTriple, t *testing.T) {
	nodes := []*GameState{node}
	for i := range movesTests {

		if !movesTests[i].preTest(nodes[i]) {
			t.Errorf("pre action test function  #%v did not pass", i)
		}

		child := nodes[i].CurrentActor().Act(nodes[i], movesTests[i].move)
		nodes = append(nodes, child)

		if !movesTests[i].postTest(child) {
			t.Errorf("post action test function  #%v did not pass", i)
		}
	}
}

func testGamePlayAfterAllMoves(node *GameState, moves []Move, test func(state *GameState) bool, t *testing.T) {
	nodes := []*GameState{node}
	for i := range moves {
		child := nodes[i].CurrentActor().Act(nodes[i], moves[i])
		nodes = append(nodes, child)
	}
	if !test(nodes[len(nodes)-1]) {
		t.Error("post game test function did not pass")
	}
}

func createRootForTest(playerAStack float64, playerBStack float64) *GameState {
	playerA := &Player{id: PlayerA, moves: nil, card: nil, stack: playerAStack}
	playerB := &Player{id: PlayerB, moves: nil, card: nil, stack: playerBStack}
	return CreateRoot(playerA, playerB)
}

func createRootWithPreparedDeck(playerAStack float64, playerBStack float64, deck *FullDeck) *GameState {
	playerA := &Player{id: PlayerA, moves: nil, card: nil, stack: playerAStack}
	playerB := &Player{id: PlayerB, moves: nil, card: nil, stack: playerBStack}
	chance := &Chance{id: ChanceId, deck: deck}

	actors := map[ActorId]Actor{PlayerA: playerA, PlayerB: playerB, ChanceId: chance}
	table := &Table{pot: 0, cards: []Card{}}

	return &GameState{round: Start, table: table,
		actors: actors, nextToMove: ChanceId, causingMove: NoMove}
}

func prepareDeckForTest(privateCardA, privateCardB, flopCard, turnCard Card) *FullDeck {
	d := CreateFullDeck(false)
	for i := range d.cards {
		if d.cards[i] == privateCardA {
			d.cards[0], d.cards[i] = d.cards[i], d.cards[0]
		}
		if d.cards[i] == privateCardB {
			d.cards[1], d.cards[i] = d.cards[i], d.cards[1]
		}
		if d.cards[i] == flopCard {
			d.cards[2], d.cards[i] = d.cards[i], d.cards[2]
		}
		if d.cards[i] == turnCard {
			d.cards[3], d.cards[i] = d.cards[i], d.cards[3]
		}
	}
	return d
}

func roundCheck(expectedRound Round) func(node *GameState) bool {
	return func(node *GameState) bool { return node.round == expectedRound }
}

func gameEnd() func(state *GameState) bool {
	return func(state *GameState) bool { return state.IsTerminal() }
}

func gameResult(result float64) func(state *GameState) bool {
	return func(state *GameState) bool {
		evaluation := state.Evaluate()
		return evaluation == result
	}
}

func noRaiseAvailable() func(state *GameState) bool {
	return func(state *GameState) bool {
		moves := state.CurrentActor().GetAvailableMoves(state)
		for _, m := range moves {
			if m == Raise {
				return false
			}
		}
		return true
	}
}

func actorToMove(actorId ActorId) func(state *GameState) bool {
	return func(state *GameState) bool {
		return state.nextToMove == actorId
	}
}

func stackEqualTo(player ActorId, stack float64) func(state *GameState) bool {
	return func(state *GameState) bool {
		return math.Abs(state.actors[player].(*Player).stack-stack) < 1e-9
	}
}

func potEqualsTo(pot float64) func(state *GameState) bool {
	return func(state *GameState) bool {
		return math.Abs(state.table.pot-pot) < 1e-9
	}
}

func noTest() func(state *GameState) bool {
	return func(state *GameState) bool {
		return true
	}
}

func onlyCheckAvailable() func(state *GameState) bool {
	return func(state *GameState) bool {
		moves := state.CurrentActor().GetAvailableMoves(state)
		if len(moves) == 1 && moves[0] == Check {
			return true
		}
		return false
	}
}
func checkAndBetAvailable() func(state *GameState) bool {
	return func(state *GameState) bool {
		moves := state.CurrentActor().GetAvailableMoves(state)
		if len(moves) == 2 && moves[0] == Check && moves[1] == Bet {
			return true
		}
		return false
	}
}

func privateCards(playerACard Card, playerBCard Card) func(state *GameState) bool {
	return func(state *GameState) bool {
		return *(state.actors[PlayerA].(*Player).card) == playerACard && *(state.actors[PlayerB].(*Player).card) == playerBCard
	}
}

func flopCard(publicFlopCard Card) func(state *GameState) bool {
	return func(state *GameState) bool {
		return state.table.cards[0] == publicFlopCard
	}
}

func turnCard(publicTurnCard Card) func(state *GameState) bool {
	return func(state *GameState) bool {
		return state.table.cards[1] == publicTurnCard
	}
}
