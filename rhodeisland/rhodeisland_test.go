package rhodeisland

import (
	. "github.com/int8/gopoker"
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
	if root.causingAction != nil {
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

	actions := root.Actions()

	if len(actions) != 52*51 {
		t.Errorf("Game root should have %v actions available, %v actions available", 52*51, len(actions))
	}

}

func TestIfParentsCorrect(t *testing.T) {
	root := createRootForTest(100., 100.)
	child := root.Act(DealPrivateCardsAction{&AceHearts, &KingClubs})
	if child.Parent() != root {
		t.Error("Root child should have root as a parent")
	}
}

func TestIfStackLimitsAvailableActions(t *testing.T) {
	root5 := createRootForTest(5., 5.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
		{DealPublicCardAction{&QueenClubs}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
		{DealPublicCardAction{&JackClubs}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
	}
	testGamePlayAfterEveryAction(root5, actionsTestsPairs, t)

	root15 := createRootForTest(15., 15.)
	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, noTest(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{DealPublicCardAction{&QueenClubs}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()}, // at this point bet size exceeds players stack
		{CheckAction, onlyCheckAvailable(), noTest()}, // only check available
		{DealPublicCardAction{&JackHearts}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()}, // at this point bet size exceeds players stack
		{CheckAction, onlyCheckAvailable(), noTest()}, // only check available
	}

	testGamePlayAfterEveryAction(root15, actionsTestsPairs, t)

	root1000Vs15 := createRootForTest(1000., 15.)
	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, noTest(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{DealPublicCardAction{&QueenClubs}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()}, // at this point bet size exceeds one of the players stack
		{CheckAction, onlyCheckAvailable(), noTest()}, // only check available
		{DealPublicCardAction{&JackHearts}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()}, // at this point bet size exceeds one of the players stack
		{CheckAction, onlyCheckAvailable(), noTest()}, // only check available
	}

	testGamePlayAfterEveryAction(root1000Vs15, actionsTestsPairs, t)

	root1000Vs35 := createRootForTest(1000., 35.)
	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, noTest(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{DealPublicCardAction{&QueenClubs}, noTest(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{DealPublicCardAction{&JackHearts}, noTest(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
	}

	testGamePlayAfterEveryAction(root1000Vs35, actionsTestsPairs, t)
}

func TestGamePlayAssertRounds(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, roundCheck(Start), roundCheck(PreFlop)},
		{CheckAction, roundCheck(PreFlop), roundCheck(PreFlop)},
		{CheckAction, roundCheck(PreFlop), roundCheck(PreFlop)},
		{DealPublicCardAction{&QueenClubs}, roundCheck(PreFlop), roundCheck(Flop)},
		{BetAction, roundCheck(Flop), roundCheck(Flop)},
		{CallAction, roundCheck(Flop), roundCheck(Flop)},
		{DealPublicCardAction{&JackClubs}, roundCheck(Flop), roundCheck(Turn)},
		{CheckAction, roundCheck(Turn), roundCheck(Turn)},
		{BetAction, roundCheck(Turn), roundCheck(Turn)},
		{CallAction, roundCheck(Turn), gameEnd()},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_MaxRaises(t *testing.T) {
	root := createRootForTest(1000., 1000.)

	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, roundCheck(Start), roundCheck(PreFlop)},
		{BetAction, roundCheck(PreFlop), roundCheck(PreFlop)},
		{RaiseAction, roundCheck(PreFlop), roundCheck(PreFlop)},
		{RaiseAction, roundCheck(PreFlop), roundCheck(PreFlop)},
		{RaiseAction, roundCheck(PreFlop), noRaiseAvailable()},
		{FoldAction, roundCheck(PreFlop), gameEnd()},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, roundCheck(Start), roundCheck(PreFlop)},
		{BetAction, roundCheck(PreFlop), roundCheck(PreFlop)},
		{RaiseAction, roundCheck(PreFlop), roundCheck(PreFlop)},
		{RaiseAction, roundCheck(PreFlop), roundCheck(PreFlop)},
		{RaiseAction, roundCheck(PreFlop), noRaiseAvailable()},
		{CallAction, roundCheck(PreFlop), roundCheck(PreFlop)},
		{DealPublicCardAction{&QueenClubs}, roundCheck(PreFlop), roundCheck(Flop)},
		{BetAction, roundCheck(Flop), roundCheck(Flop)},
		{RaiseAction, roundCheck(Flop), roundCheck(Flop)},
		{RaiseAction, roundCheck(Flop), roundCheck(Flop)},
		{RaiseAction, roundCheck(Flop), noRaiseAvailable()},
		{CallAction, roundCheck(Flop), roundCheck(Flop)},
		{DealPublicCardAction{&JackClubs}, roundCheck(Flop), roundCheck(Turn)},
		{CheckAction, roundCheck(Turn), roundCheck(Turn)},
		{CheckAction, roundCheck(Turn), gameEnd()},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

}

func TestGamePlay_CheckIfPlayerToMoveCorrect(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, actorToMove(ChanceId), actorToMove(PlayerA)},
		{CheckAction, actorToMove(PlayerA), actorToMove(PlayerB)},
		{CheckAction, actorToMove(PlayerB), actorToMove(ChanceId)},
		{DealPublicCardAction{&QueenClubs}, actorToMove(ChanceId), actorToMove(PlayerA)},
		{BetAction, actorToMove(PlayerA), actorToMove(PlayerB)},
		{CallAction, actorToMove(PlayerB), actorToMove(ChanceId)},
		{DealPublicCardAction{&JackHearts}, actorToMove(ChanceId), actorToMove(PlayerA)},
		{CheckAction, actorToMove(PlayerA), actorToMove(PlayerB)},
		{BetAction, actorToMove(PlayerB), actorToMove(PlayerA)},
		{CallAction, actorToMove(PlayerA), noTest()},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfStacksChange(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, stackEqualsTo(PlayerA, 100.), stackEqualsTo(PlayerA, 100.-Ante)},
		{CheckAction, stackEqualsTo(PlayerA, 100.-Ante), stackEqualsTo(PlayerA, 100.-Ante)},
		{CheckAction, stackEqualsTo(PlayerB, 100.-Ante), stackEqualsTo(PlayerB, 100.-Ante)},
		{DealPublicCardAction{&JackHearts}, stackEqualsTo(PlayerA, 100.-Ante), stackEqualsTo(PlayerA, 100.-Ante)},
		{BetAction, stackEqualsTo(PlayerA, 100.-Ante), stackEqualsTo(PlayerA, 100.-Ante-PostFlopBetSize)},
		{CallAction, stackEqualsTo(PlayerB, 100.-Ante), stackEqualsTo(PlayerB, 100.-Ante-PostFlopBetSize)},
		{DealPublicCardAction{&QueenSpades}, stackEqualsTo(PlayerB, 100.-Ante-PostFlopBetSize), stackEqualsTo(PlayerB, 100.-Ante-PostFlopBetSize)},
		{CheckAction, stackEqualsTo(PlayerA, 100.-Ante-PostFlopBetSize), stackEqualsTo(PlayerA, 100.-Ante-PostFlopBetSize)},
		{BetAction, stackEqualsTo(PlayerB, 100.-Ante-PostFlopBetSize), stackEqualsTo(PlayerB, 100.-Ante-2*PostFlopBetSize)},
		{CallAction, stackEqualsTo(PlayerA, 100.-Ante-PostFlopBetSize), stackEqualsTo(PlayerA, 100.-Ante-2*PostFlopBetSize)},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfPotChanges(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, potEqualsTo(0.0), potEqualsTo(10.0)},
		{CheckAction, potEqualsTo(10), potEqualsTo(10)},
		{CheckAction, potEqualsTo(10), potEqualsTo(10)},
		{DealPublicCardAction{&JackHearts}, potEqualsTo(10), potEqualsTo(10)},
		{BetAction, potEqualsTo(10), potEqualsTo(10 + PostFlopBetSize)},
		{CallAction, potEqualsTo(10 + PostFlopBetSize), potEqualsTo(10 + 2*PostFlopBetSize)},
		{DealPublicCardAction{&QueenSpades}, potEqualsTo(10 + 2*PostFlopBetSize), potEqualsTo(10 + 2*PostFlopBetSize)},
		{CheckAction, potEqualsTo(10 + 2*PostFlopBetSize), potEqualsTo(10 + 2*PostFlopBetSize)},
		{BetAction, potEqualsTo(10 + 2*PostFlopBetSize), potEqualsTo(10 + 3*PostFlopBetSize)},
		{CallAction, potEqualsTo(10 + 3*PostFlopBetSize), potEqualsTo(10 + 4*PostFlopBetSize)},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestIfCardsGoToPlayers(t *testing.T) {

	root := createRootForTest(100., 100)

	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, noTest(), privateCards(AceHearts, KingClubs)},
		{CheckAction, noTest(), noTest()},
		{CheckAction, noTest(), noTest()},
		{DealPublicCardAction{&JackHearts}, noTest(), flopCard(JackHearts)},
		{CheckAction, noTest(), noTest()},
		{CheckAction, noTest(), noTest()},
		{DealPublicCardAction{&QueenSpades}, noTest(), turnCard(QueenSpades)},
		{CheckAction, noTest(), noTest()},
		{CheckAction, noTest(), noTest()},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfChildPointersDifferFromParentsPointers(t *testing.T) {
	root := createRootForTest(100., 100.)
	child := root.Act(DealPrivateCardsAction{&AceHearts, &KingClubs})

	if child.(*RIGameState).actors[ChanceId] == root.actors[ChanceId] {
		t.Error("chance actor refers to the same value in both child and parent")
	}

	if child.(*RIGameState).actors[PlayerA] == root.actors[PlayerA] {
		t.Error("PlayerA actor refers to the same value in both child and parent")
	}

	if child.(*RIGameState).actors[PlayerB] == root.actors[PlayerB] {
		t.Error("PlayerB actor refers to the same value in both child and parent")
	}

	if child.(*RIGameState).table == root.table {
		t.Error("table should be different for child and parent")
	}
}

func TestGamePlayEvaluationFlushVsNothing(t *testing.T) {
	root := createRootForTest(100., 100.)
	hands := DealPrivateCardsAction{&AceHearts, &C2Spades}
	flop := DealPublicCardAction{&JackHearts}
	turn := DealPublicCardAction{&KingHearts}
	actions := []Action{hands, CheckAction, CheckAction, flop, BetAction, CallAction, turn, CheckAction, BetAction, CallAction}

	singlePlayerPotContribution := Ante + PostFlopBetSize + PostFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationFlushVsStraightFlush(t *testing.T) {
	hands := DealPrivateCardsAction{&AceHearts, &QueenHearts}
	flop := DealPublicCardAction{&JackHearts}
	turn := DealPublicCardAction{&KingHearts}

	root := createRootForTest(100., 100.)

	actions := []Action{hands, BetAction, CallAction, flop, BetAction, CallAction, turn, CheckAction, BetAction, CallAction}
	singlePlayerPotContribution := Ante + PreFlopBetSize + PostFlopBetSize + PostFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationPairVsPairDraw(t *testing.T) {
	hands := DealPrivateCardsAction{&AceHearts, &AceSpades}
	flop := DealPublicCardAction{&JackSpades}
	turn := DealPublicCardAction{&KingHearts}

	root := createRootForTest(100., 100.)

	actions := []Action{hands, CheckAction, CheckAction, flop, CheckAction, CheckAction, turn, CheckAction, CheckAction}
	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(0), t)
}

func TestGamePlayEvaluationPairVsPairAWins(t *testing.T) {
	hands := DealPrivateCardsAction{&AceHearts, &KingSpades}
	flop := DealPublicCardAction{&AceSpades}
	turn := DealPublicCardAction{&KingHearts}

	root := createRootForTest(100., 100.)

	actions := []Action{hands, CheckAction, CheckAction, flop, CheckAction, CheckAction, turn, CheckAction, CheckAction}
	singlePlayerPotContribution := Ante
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationPairVsPairBWinsBetterOwnCard(t *testing.T) {
	hands := DealPrivateCardsAction{&JackHearts, &KingSpades}
	flop := DealPublicCardAction{&C2Spades}
	turn := DealPublicCardAction{&C2Hearts}

	root := createRootForTest(100., 100.)

	actions := []Action{hands, CheckAction, CheckAction, flop, CheckAction, CheckAction, turn, CheckAction, CheckAction}
	singlePlayerPotContribution := Ante
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationStraightVsStraightAWinsBetterOwnCard(t *testing.T) {

	hands := DealPrivateCardsAction{&KingHearts, &C10Spades}
	flop := DealPublicCardAction{&JackClubs}
	turn := DealPublicCardAction{&QueenDiamonds}

	root := createRootForTest(100., 100.)

	actions := []Action{hands, BetAction, CallAction, flop, CheckAction, CheckAction, turn, CheckAction, CheckAction}
	singlePlayerPotContribution := Ante + PreFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationPairVsThreeOfAKindBWins(t *testing.T) {

	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	flop := DealPublicCardAction{&KingDiamonds}
	turn := DealPublicCardAction{&KingSpades}

	root := createRootForTest(100., 100.)

	actions := []Action{hands, BetAction, CallAction, flop, CheckAction, CheckAction, turn, CheckAction, CheckAction}

	singlePlayerPotContribution := Ante + PreFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationOwnCardVsOwnCardAWins(t *testing.T) {

	hands := DealPrivateCardsAction{&KingHearts, &C10Clubs}
	flop := DealPublicCardAction{&C2Diamonds}
	turn := DealPublicCardAction{&C7Spades}

	root := createRootForTest(100., 100.)
	actions := []Action{hands, BetAction, CallAction, flop, CheckAction, CheckAction, turn, BetAction, CallAction}

	singlePlayerPotContribution := Ante + PreFlopBetSize + PostFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAFoldsOnTurn(t *testing.T) {

	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	flop := DealPublicCardAction{&KingDiamonds}
	turn := DealPublicCardAction{&KingSpades}

	root := createRootForTest(100., 100.)
	actions := []Action{hands, BetAction, CallAction, flop, CheckAction, CheckAction, turn, CheckAction, BetAction, FoldAction}

	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsOnTurn(t *testing.T) {
	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	flop := DealPublicCardAction{&KingDiamonds}
	turn := DealPublicCardAction{&KingSpades}

	root := createRootForTest(100., 100.)

	actions := []Action{hands, BetAction, CallAction, flop, CheckAction, CheckAction, turn, BetAction, FoldAction}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAFoldsOnFlop(t *testing.T) {
	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	flop := DealPublicCardAction{&KingDiamonds}

	root := createRootForTest(100., 100.)
	actions := []Action{hands, BetAction, CallAction, flop, CheckAction, BetAction, FoldAction}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsOnFlop(t *testing.T) {
	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	flop := DealPublicCardAction{&KingDiamonds}

	root := createRootForTest(100., 100.)
	actions := []Action{hands, BetAction, CallAction, flop, BetAction, FoldAction}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAFoldsPreFlop(t *testing.T) {
	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, CheckAction, BetAction, FoldAction}

	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsPreFlop(t *testing.T) {
	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, CheckAction, BetAction, RaiseAction, FoldAction}

	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsPreFlopManyRaises(t *testing.T) {
	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, CheckAction, BetAction, RaiseAction, RaiseAction, RaiseAction, FoldAction}

	singlePlayerPotContribution := Ante + 3*PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayInformationSetForAAfterRaises(t *testing.T) {
	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, CheckAction, BetAction, RaiseAction, RaiseAction}

	targetInformationSet := createInformationSet(C10Hearts, NoCard, NoCard, actions)
	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForB_ChecksOnly(t *testing.T) {

	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	flop := DealPublicCardAction{&KingDiamonds}
	turn := DealPublicCardAction{&KingSpades}
	root := createRootForTest(100., 100.)

	actions := []Action{hands, CheckAction, CheckAction, flop, CheckAction, CheckAction, turn, CheckAction}
	targetInformationSet := createInformationSet(KingClubs, KingDiamonds, KingSpades, actions)
	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForA_NoActions(t *testing.T) {
	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}

	root := createRootForTest(100., 100.)

	actions := []Action{hands}
	targetInformationSet := createInformationSet(C10Hearts, NoCard, NoCard, actions)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForB_SingleCheck(t *testing.T) {

	hands := DealPrivateCardsAction{&C10Hearts, &AceSpades}
	root := createRootForTest(100., 100.)

	actions := []Action{hands, CheckAction}
	targetInformationSet := createInformationSet(AceSpades, NoCard, NoCard, actions)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForB_BetCallAndChecksOnly(t *testing.T) {

	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	flop := DealPublicCardAction{&KingDiamonds}
	turn := DealPublicCardAction{&KingSpades}
	root := createRootForTest(100., 100.)

	actions := []Action{hands, CheckAction, BetAction, CallAction, flop, CheckAction, CheckAction, turn, CheckAction}
	targetInformationSet := createInformationSet(KingClubs, KingDiamonds, KingSpades, actions)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForBAfterCheckBetRaise(t *testing.T) {

	hands := DealPrivateCardsAction{&C10Hearts, &QueenClubs}
	root := createRootForTest(100., 100.)

	actions := []Action{hands, CheckAction, BetAction, RaiseAction}
	targetInformationSet := createInformationSet(QueenClubs, NoCard, NoCard, actions)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}
func testGamePlayAfterEveryAction(node *RIGameState, actionsTests []ActionTestsTriple, t *testing.T) {
	nodes := []GameState{node}
	for i := range actionsTests {

		if !actionsTests[i].preTest(nodes[i].(*RIGameState)) {
			t.Errorf("pre action test function  #%v did not pass", i)
		}

		child := nodes[i].Act(actionsTests[i].Action)
		nodes = append(nodes, child)

		if !actionsTests[i].postTest(child.(*RIGameState)) {
			t.Errorf("post action test function  #%v did not pass", i)
		}
	}
}

func testGamePlayAfterAllActions(node *RIGameState, actions []Action, test func(state *RIGameState) bool, t *testing.T) {
	nodes := []GameState{node}
	for i := range actions {
		child := nodes[i].Act(actions[i])
		nodes = append(nodes, child)
	}
	if !test(nodes[len(nodes)-1].(*RIGameState)) {
		t.Error("post game test function did not pass")
	}
}

func createRootForTest(playerAStack float32, playerBStack float32) *RIGameState {
	playerA := &Player{Id: PlayerA, Actions: nil, Card: nil, Stack: playerAStack}
	playerB := &Player{Id: PlayerB, Actions: nil, Card: nil, Stack: playerBStack}
	return Root(playerA, playerB, CreateFullDeck(true))
}

func roundCheck(expectedRound Round) func(node *RIGameState) bool {
	return func(node *RIGameState) bool { return node.round == expectedRound }
}

func gameEnd() func(state *RIGameState) bool {
	return func(state *RIGameState) bool { return state.IsTerminal() }
}

func gameResult(result float32) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		evaluation := state.Evaluate()
		return evaluation == result
	}
}

func noRaiseAvailable() func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		Actions := state.Actions()
		for _, m := range Actions {
			if m.Name() == Raise {
				return false
			}
		}
		return true
	}
}

func actorToMove(actorId ActorId) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return state.nextToMove == actorId
	}
}

func stackEqualsTo(player ActorId, stack float32) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return math.Abs(float64(state.actors[player].(*Player).Stack-stack)) < 1e-9
	}
}

func potEqualsTo(pot float32) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return math.Abs(float64(state.table.Pot-pot)) < 1e-9
	}
}

func noTest() func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return true
	}
}

func onlyCheckAvailable() func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		Actions := state.Actions()
		if len(Actions) == 1 && Actions[0].Name() == Check {
			return true
		}
		return false
	}
}
func checkAndBetAvailable() func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		Actions := state.Actions()
		if len(Actions) == 2 && Actions[0].Name() == Check && Actions[1].Name() == Bet {
			return true
		}
		return false
	}
}

func privateCards(playerACard Card, playerBCard Card) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return *(state.actors[PlayerA].(*Player).Card) == playerACard && *(state.actors[PlayerB].(*Player).Card) == playerBCard
	}
}

func flopCard(publicFlopCard Card) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return state.table.Cards[0] == publicFlopCard
	}
}

func turnCard(publicTurnCard Card) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return state.table.Cards[1] == publicTurnCard
	}
}

func lastInformationSet(informationSet [InformationSetSize]bool) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		currentInformationSet := state.InformationSet()
		return currentInformationSet == informationSet
	}
}

func createInformationSet(prvCard Card, flopCard Card, turnCard Card, actions []Action) [InformationSetSize]bool {

	informationSet := [InformationSetSize]bool{
		prvCard.Symbol[0], prvCard.Symbol[1], prvCard.Symbol[2], prvCard.Symbol[3],
		prvCard.Suit[0], prvCard.Suit[1], prvCard.Suit[2],
		flopCard.Symbol[0], flopCard.Symbol[1], flopCard.Symbol[2], flopCard.Symbol[3],
		flopCard.Suit[0], flopCard.Suit[1], flopCard.Suit[2],
		turnCard.Symbol[0], turnCard.Symbol[1], turnCard.Symbol[2], turnCard.Symbol[3],
		turnCard.Suit[0], turnCard.Suit[1], turnCard.Suit[2],
	}

	var currentAction Action
	for i := 21; len(actions) > 0; i += 3 {
		// somehow tricky pop..
		currentAction, actions = actions[len(actions)-1], actions[:len(actions)-1]
		informationSet[i] = currentAction.Name()[0]
		informationSet[i+1] = currentAction.Name()[1]
		informationSet[i+2] = currentAction.Name()[2]
	}

	return informationSet
}
