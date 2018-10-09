package gocfr

import (
	"math"
	"testing"
)

type ActionTestsTriple struct {
	Action   Action
	preTest  func(state *RIGameState) bool
	postTest func(state *RIGameState) bool
}

func TestGameCreation(t *testing.T) {
	root := createRootForTest(100., 100.)
	if root.causingAction != NoAction {
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

	Actions := root.actors[root.nextToMove].GetAvailableActions(root)

	if Actions == nil {
		t.Error("Game root should have one action available, no actions available")
	}

	if len(Actions) != 1 {
		t.Errorf("Game root should have one action available, %v actions available", len(Actions))
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
	actionsTestsPairs := []ActionTestsTriple{
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
	testGamePlayAfterEveryAction(root5, actionsTestsPairs, t)

	root15 := createRootForTest(15., 15.)
	actionsTestsPairs = []ActionTestsTriple{
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

	testGamePlayAfterEveryAction(root15, actionsTestsPairs, t)

	root1000Vs15 := createRootForTest(1000., 15.)
	actionsTestsPairs = []ActionTestsTriple{
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

	testGamePlayAfterEveryAction(root1000Vs15, actionsTestsPairs, t)

	root1000Vs35 := createRootForTest(1000., 35.)
	actionsTestsPairs = []ActionTestsTriple{
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

	testGamePlayAfterEveryAction(root1000Vs35, actionsTestsPairs, t)
}

func TestGamePlay_1(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
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

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_MaxRaises(t *testing.T) {
	root := createRootForTest(1000., 1000.)

	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCards, roundCheck(Start), roundCheck(PreFlop)},
		{Bet, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), noRaiseAvailable()},
		{Fold, roundCheck(PreFlop), gameEnd()},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCards, roundCheck(Start), roundCheck(PreFlop)},
		{Bet, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), noRaiseAvailable()},
		{Call, roundCheck(PreFlop), roundCheck(PreFlop)},
		{DealPublicCard, roundCheck(PreFlop), roundCheck(Flop)},
		{Bet, roundCheck(Flop), roundCheck(Flop)},
		{Raise, roundCheck(Flop), roundCheck(Flop)},
		{Raise, roundCheck(Flop), roundCheck(Flop)},
		{Raise, roundCheck(Flop), noRaiseAvailable()},
		{Call, roundCheck(Flop), roundCheck(Flop)},
		{DealPublicCard, roundCheck(Flop), roundCheck(Turn)},
		{Check, roundCheck(Turn), roundCheck(Turn)},
		{Check, roundCheck(Turn), gameEnd()},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

}

func TestGamePlay_CheckIfPlayerToActionCorrect(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCards, actorToAction(ChanceId), actorToAction(PlayerA)},
		{Check, actorToAction(PlayerA), actorToAction(PlayerB)},
		{Check, actorToAction(PlayerB), actorToAction(ChanceId)},
		{DealPublicCard, actorToAction(ChanceId), actorToAction(PlayerA)},
		{Bet, actorToAction(PlayerA), actorToAction(PlayerB)},
		{Call, actorToAction(PlayerB), actorToAction(ChanceId)},
		{DealPublicCard, actorToAction(ChanceId), actorToAction(PlayerA)},
		{Check, actorToAction(PlayerA), actorToAction(PlayerB)},
		{Bet, actorToAction(PlayerB), actorToAction(PlayerA)},
		{Call, actorToAction(PlayerA), noTest()},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfStacksChange(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
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

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfPotChanges(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
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

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestIfRootCreationWithDeckPreparedWorks(t *testing.T) {
	aceHearts := Card{Ace, Hearts}
	c2Spades := Card{C2, Spades}
	jackHearts := Card{Jack, Hearts}
	kingHearts := Card{King, Hearts}

	preparedDeck := prepareDeckForTest(aceHearts, c2Spades, jackHearts, kingHearts)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actionsTestsPairs := []ActionTestsTriple{
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

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

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

	actions := []Action{DealPrivateCards, Check, Check, DealPublicCard, Bet, Call, DealPublicCard, Check, Bet, Call}
	singlePlayerPotContribution := Ante + PostFlopBetSize + PostFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationFlushVsStraightFlush(t *testing.T) {
	privateACard := Card{Ace, Hearts}
	privateBCard := Card{Queen, Hearts}
	flopPublicCard := Card{Jack, Hearts}
	turnPublicCard := Card{King, Hearts}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Bet, Call, DealPublicCard, Bet, Call, DealPublicCard, Check, Bet, Call}
	singlePlayerPotContribution := Ante + PreFlopBetSize + PostFlopBetSize + PostFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationPairVsPairDraw(t *testing.T) {
	privateACard := Card{Ace, Hearts}
	privateBCard := Card{Ace, Spades}
	flopPublicCard := Card{Jack, Spades}
	turnPublicCard := Card{King, Hearts}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Check, Check, DealPublicCard, Check, Check, DealPublicCard, Check, Check}
	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(0), t)
}

func TestGamePlayEvaluationPairVsPairAWins(t *testing.T) {
	privateACard := Card{Ace, Hearts}
	privateBCard := Card{King, Spades}
	flopPublicCard := Card{Ace, Spades}
	turnPublicCard := Card{King, Hearts}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Check, Check, DealPublicCard, Check, Check, DealPublicCard, Check, Check}
	singlePlayerPotContribution := Ante
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationPairVsPairBWinsBetterOwnCard(t *testing.T) {
	privateACard := Card{Jack, Hearts}
	privateBCard := Card{King, Spades}
	flopPublicCard := Card{C2, Spades}
	turnPublicCard := Card{C2, Hearts}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Check, Check, DealPublicCard, Check, Check, DealPublicCard, Check, Check}
	singlePlayerPotContribution := Ante
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationStraightVsStraightAWinsBetterOwnCard(t *testing.T) {
	privateACard := Card{King, Hearts}
	privateBCard := Card{C10, Spades}
	flopPublicCard := Card{Jack, Clubs}
	turnPublicCard := Card{Queen, Diamonds}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Bet, Call, DealPublicCard, Check, Check, DealPublicCard, Check, Check}
	singlePlayerPotContribution := Ante + PreFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationPairVsThreeOfAKindBWins(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Bet, Call, DealPublicCard, Check, Check, DealPublicCard, Check, Check}
	singlePlayerPotContribution := Ante + PreFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationOwnCardVsOwnCardAWins(t *testing.T) {
	privateACard := Card{King, Hearts}
	privateBCard := Card{C10, Clubs}
	flopPublicCard := Card{C2, Diamonds}
	turnPublicCard := Card{C7, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Bet, Call, DealPublicCard, Check, Check, DealPublicCard, Bet, Call}
	singlePlayerPotContribution := Ante + PreFlopBetSize + PostFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAFoldsOnTurn(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Bet, Call, DealPublicCard, Check, Check, DealPublicCard, Check, Bet, Fold}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsOnTurn(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Bet, Call, DealPublicCard, Check, Check, DealPublicCard, Bet, Fold}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAFoldsOnFlop(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Bet, Call, DealPublicCard, Check, Bet, Fold}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsOnFlop(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Bet, Call, DealPublicCard, Bet, Fold}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAFoldsPreFlop(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Check, Bet, Fold}
	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsPreFlop(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Check, Bet, Raise, Fold}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsPreFlopManyRaises(t *testing.T) {
	privateACard := Card{C10, Hearts}
	privateBCard := Card{King, Clubs}
	flopPublicCard := Card{King, Diamonds}
	turnPublicCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(privateACard, privateBCard, flopPublicCard, turnPublicCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Check, Bet, Raise, Raise, Raise, Fold}
	singlePlayerPotContribution := Ante + 3*PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayInformationSetForAAfterRaises(t *testing.T) {
	aCard := Card{C10, Hearts}
	bCard := Card{King, Clubs}
	flopCard := Card{King, Diamonds}
	turnCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(aCard, bCard, flopCard, turnCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Check, Bet, Raise, Raise}
	targetInformationSet := [InformationSetSize]byte{byte(aCard.name), byte(aCard.suit)}
	targetInformationSet[6] = byte(Raise)
	targetInformationSet[7] = byte(Raise)
	targetInformationSet[8] = byte(Bet)
	targetInformationSet[9] = byte(Check)
	targetInformationSet[10] = byte(DealPrivateCards)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForB_ChecksOnly(t *testing.T) {
	aCard := Card{C10, Hearts}
	bCard := Card{King, Clubs}
	flopCard := Card{King, Diamonds}
	turnCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(aCard, bCard, flopCard, turnCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Check, Check, DealPublicCard, Check, Check, DealPublicCard, Check}
	targetInformationSet := [InformationSetSize]byte{byte(bCard.name), byte(bCard.suit), byte(flopCard.name),
		byte(flopCard.suit), byte(turnCard.name), byte(turnCard.suit)}
	targetInformationSet[6] = byte(Check)
	targetInformationSet[7] = byte(DealPublicCard)
	targetInformationSet[8] = byte(Check)
	targetInformationSet[9] = byte(Check)
	targetInformationSet[10] = byte(DealPublicCard)
	targetInformationSet[11] = byte(Check)
	targetInformationSet[12] = byte(Check)
	targetInformationSet[13] = byte(DealPrivateCards)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForA_NoActions(t *testing.T) {
	aCard := Card{C10, Hearts}
	bCard := Card{King, Clubs}
	flopCard := Card{King, Diamonds}
	turnCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(aCard, bCard, flopCard, turnCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards}
	targetInformationSet := [InformationSetSize]byte{byte(aCard.name), byte(aCard.suit)}
	targetInformationSet[6] = byte(DealPrivateCards)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForB_SingleCheck(t *testing.T) {
	aCard := Card{C10, Hearts}
	bCard := Card{Ace, Spades}
	flopCard := Card{King, Diamonds}
	turnCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(aCard, bCard, flopCard, turnCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Check}
	targetInformationSet := [InformationSetSize]byte{byte(bCard.name), byte(bCard.suit)}
	targetInformationSet[6] = byte(Check)
	targetInformationSet[7] = byte(DealPrivateCards)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForB_BetCallAndChecksOnly(t *testing.T) {
	aCard := Card{C10, Hearts}
	bCard := Card{King, Clubs}
	flopCard := Card{King, Diamonds}
	turnCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(aCard, bCard, flopCard, turnCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Check, Bet, Call, DealPublicCard, Check, Check, DealPublicCard, Check}
	targetInformationSet := [InformationSetSize]byte{byte(bCard.name), byte(bCard.suit), byte(flopCard.name),
		byte(flopCard.suit), byte(turnCard.name), byte(turnCard.suit)}
	targetInformationSet[6] = byte(Check)
	targetInformationSet[7] = byte(DealPublicCard)
	targetInformationSet[8] = byte(Check)
	targetInformationSet[9] = byte(Check)
	targetInformationSet[10] = byte(DealPublicCard)
	targetInformationSet[11] = byte(Call)
	targetInformationSet[12] = byte(Bet)
	targetInformationSet[13] = byte(Check)
	targetInformationSet[14] = byte(DealPrivateCards)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForBAfterCheckBetRaise(t *testing.T) {
	aCard := Card{C10, Hearts}
	bCard := Card{Queen, Clubs}
	flopCard := Card{King, Diamonds}
	turnCard := Card{King, Spades}

	preparedDeck := prepareDeckForTest(aCard, bCard, flopCard, turnCard)
	root := createRootWithPreparedDeck(100., 100., preparedDeck)

	actions := []Action{DealPrivateCards, Check, Bet, Raise}
	targetInformationSet := [InformationSetSize]byte{byte(bCard.name), byte(bCard.suit)}
	targetInformationSet[6] = byte(Raise)
	targetInformationSet[7] = byte(Bet)
	targetInformationSet[8] = byte(Check)
	targetInformationSet[9] = byte(DealPrivateCards)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func testGamePlayAfterEveryAction(node *RIGameState, actionsTests []ActionTestsTriple, t *testing.T) {
	nodes := []GameStateHolder{node}
	for i := range actionsTests {

		if !actionsTests[i].preTest(nodes[i].(*RIGameState)) {
			t.Errorf("pre action test function  #%v did not pass", i)
		}

		child := nodes[i].Child(actionsTests[i].Action)
		nodes = append(nodes, child)

		if !actionsTests[i].postTest(child.(*RIGameState)) {
			t.Errorf("post action test function  #%v did not pass", i)
		}
	}
}

func testGamePlayAfterAllActions(node *RIGameState, actions []Action, test func(state *RIGameState) bool, t *testing.T) {
	nodes := []GameStateHolder{node}
	for i := range actions {
		child := nodes[i].Child(actions[i])
		nodes = append(nodes, child)
	}
	if !test(nodes[len(nodes)-1].(*RIGameState)) {
		t.Error("post game test function did not pass")
	}
}

func createRootForTest(playerAStack float64, playerBStack float64) *RIGameState {
	playerA := &Player{id: PlayerA, actions: nil, card: nil, stack: playerAStack}
	playerB := &Player{id: PlayerB, actions: nil, card: nil, stack: playerBStack}
	return CreateRoot(playerA, playerB)
}

func createRootWithPreparedDeck(playerAStack float64, playerBStack float64, deck *FullDeck) *RIGameState {
	playerA := &Player{id: PlayerA, actions: nil, card: nil, stack: playerAStack}
	playerB := &Player{id: PlayerB, actions: nil, card: nil, stack: playerBStack}
	chance := &Chance{id: ChanceId, deck: deck}

	actors := map[ActorId]Actor{PlayerA: playerA, PlayerB: playerB, ChanceId: chance}
	table := &Table{pot: 0, cards: []Card{}}

	return &RIGameState{round: Start, table: table,
		actors: actors, parent: nil, nextToMove: ChanceId, causingAction: NoAction}
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

func roundCheck(expectedRound Round) func(node *RIGameState) bool {
	return func(node *RIGameState) bool { return node.round == expectedRound }
}

func gameEnd() func(state *RIGameState) bool {
	return func(state *RIGameState) bool { return state.IsTerminal() }
}

func gameResult(result float64) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		evaluation := state.Evaluate()
		return evaluation == result
	}
}

func noRaiseAvailable() func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		Actions := state.CurrentActor().GetAvailableActions(state)
		for _, m := range Actions {
			if m == Raise {
				return false
			}
		}
		return true
	}
}

func actorToAction(actorId ActorId) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return state.nextToMove == actorId
	}
}

func stackEqualTo(player ActorId, stack float64) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return math.Abs(state.actors[player].(*Player).stack-stack) < 1e-9
	}
}

func potEqualsTo(pot float64) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return math.Abs(state.table.pot-pot) < 1e-9
	}
}

func noTest() func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return true
	}
}

func onlyCheckAvailable() func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		Actions := state.CurrentActor().GetAvailableActions(state)
		if len(Actions) == 1 && Actions[0] == Check {
			return true
		}
		return false
	}
}
func checkAndBetAvailable() func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		Actions := state.CurrentActor().GetAvailableActions(state)
		if len(Actions) == 2 && Actions[0] == Check && Actions[1] == Bet {
			return true
		}
		return false
	}
}

func privateCards(playerACard Card, playerBCard Card) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return *(state.actors[PlayerA].(*Player).card) == playerACard && *(state.actors[PlayerB].(*Player).card) == playerBCard
	}
}

func flopCard(publicFlopCard Card) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return state.table.cards[0] == publicFlopCard
	}
}

func turnCard(publicTurnCard Card) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return state.table.cards[1] == publicTurnCard
	}
}

func lastInformationSet(informationSet [InformationSetSize]byte) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		currentInformationSet := state.CurrentInformationSet()
		return currentInformationSet == informationSet
	}
}
