package rhodeisland

import (
	"github.com/int8/go-counterfactual-regret-minimization/acting"
	"github.com/int8/go-counterfactual-regret-minimization/cards"
	"github.com/int8/go-counterfactual-regret-minimization/games"
	"github.com/int8/go-counterfactual-regret-minimization/rounds"
	"math"
	"testing"
)

type ActionTestsTriple struct {
	Action   acting.Action
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

	if root.round != rounds.Start {
		t.Error("Initial round of the game should be rounds.Start")
	}

	if root.IsTerminal() == true {
		t.Error("Game root should not be terminal")
	}

	actions := root.Actions()

	if len(actions) != 52*51 {
		t.Errorf("Game root should have %v acting available, %v acting available", 52*51, len(actions))
	}

}

func TestIfParentsCorrect(t *testing.T) {
	root := createRootForTest(100., 100.)
	child := root.Act(DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs})
	if child.Parent() != root {
		t.Error("Root child should have root as a parent")
	}
}

func TestIfStackLimitsAvailableActions(t *testing.T) {
	root5 := createRootForTest(5., 5.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
		{DealPublicCardAction{&cards.QueenClubs}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
		{DealPublicCardAction{&cards.JackClubs}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
	}
	testGamePlayAfterEveryAction(root5, actionsTestsPairs, t)

	root15 := createRootForTest(15., 15.)
	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs}, noTest(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{DealPublicCardAction{&cards.QueenClubs}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()}, // at this point bet size exceeds players stack
		{CheckAction, onlyCheckAvailable(), noTest()}, // only check available
		{DealPublicCardAction{&cards.JackHearts}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()}, // at this point bet size exceeds players stack
		{CheckAction, onlyCheckAvailable(), noTest()}, // only check available
	}

	testGamePlayAfterEveryAction(root15, actionsTestsPairs, t)

	root1000Vs15 := createRootForTest(1000., 15.)
	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs}, noTest(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{DealPublicCardAction{&cards.QueenClubs}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()}, // at this point bet size exceeds one of the players stack
		{CheckAction, onlyCheckAvailable(), noTest()}, // only check available
		{DealPublicCardAction{&cards.JackHearts}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()}, // at this point bet size exceeds one of the players stack
		{CheckAction, onlyCheckAvailable(), noTest()}, // only check available
	}

	testGamePlayAfterEveryAction(root1000Vs15, actionsTestsPairs, t)

	root1000Vs35 := createRootForTest(1000., 35.)
	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs}, noTest(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{DealPublicCardAction{&cards.QueenClubs}, noTest(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{DealPublicCardAction{&cards.JackHearts}, noTest(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
	}

	testGamePlayAfterEveryAction(root1000Vs35, actionsTestsPairs, t)
}

func TestGamePlayAssertRounds(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs}, roundCheck(rounds.Start), roundCheck(rounds.PreFlop)},
		{CheckAction, roundCheck(rounds.PreFlop), roundCheck(rounds.PreFlop)},
		{CheckAction, roundCheck(rounds.PreFlop), roundCheck(rounds.PreFlop)},
		{DealPublicCardAction{&cards.QueenClubs}, roundCheck(rounds.PreFlop), roundCheck(rounds.Flop)},
		{BetAction, roundCheck(rounds.Flop), roundCheck(rounds.Flop)},
		{CallAction, roundCheck(rounds.Flop), roundCheck(rounds.Flop)},
		{DealPublicCardAction{&cards.JackClubs}, roundCheck(rounds.Flop), roundCheck(rounds.Turn)},
		{CheckAction, roundCheck(rounds.Turn), roundCheck(rounds.Turn)},
		{BetAction, roundCheck(rounds.Turn), roundCheck(rounds.Turn)},
		{CallAction, roundCheck(rounds.Turn), gameEnd()},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_MaxRaises(t *testing.T) {
	root := createRootForTest(1000., 1000.)

	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs}, roundCheck(rounds.Start), roundCheck(rounds.PreFlop)},
		{BetAction, roundCheck(rounds.PreFlop), roundCheck(rounds.PreFlop)},
		{RaiseAction, roundCheck(rounds.PreFlop), roundCheck(rounds.PreFlop)},
		{RaiseAction, roundCheck(rounds.PreFlop), roundCheck(rounds.PreFlop)},
		{RaiseAction, roundCheck(rounds.PreFlop), noRaiseAvailable()},
		{FoldAction, roundCheck(rounds.PreFlop), gameEnd()},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs}, roundCheck(rounds.Start), roundCheck(rounds.PreFlop)},
		{BetAction, roundCheck(rounds.PreFlop), roundCheck(rounds.PreFlop)},
		{RaiseAction, roundCheck(rounds.PreFlop), roundCheck(rounds.PreFlop)},
		{RaiseAction, roundCheck(rounds.PreFlop), roundCheck(rounds.PreFlop)},
		{RaiseAction, roundCheck(rounds.PreFlop), noRaiseAvailable()},
		{CallAction, roundCheck(rounds.PreFlop), roundCheck(rounds.PreFlop)},
		{DealPublicCardAction{&cards.QueenClubs}, roundCheck(rounds.PreFlop), roundCheck(rounds.Flop)},
		{BetAction, roundCheck(rounds.Flop), roundCheck(rounds.Flop)},
		{RaiseAction, roundCheck(rounds.Flop), roundCheck(rounds.Flop)},
		{RaiseAction, roundCheck(rounds.Flop), roundCheck(rounds.Flop)},
		{RaiseAction, roundCheck(rounds.Flop), noRaiseAvailable()},
		{CallAction, roundCheck(rounds.Flop), roundCheck(rounds.Flop)},
		{DealPublicCardAction{&cards.JackClubs}, roundCheck(rounds.Flop), roundCheck(rounds.Turn)},
		{CheckAction, roundCheck(rounds.Turn), roundCheck(rounds.Turn)},
		{CheckAction, roundCheck(rounds.Turn), gameEnd()},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

}

func TestGamePlay_CheckIfPlayerToMoveCorrect(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs}, actorToMove(acting.ChanceId), actorToMove(acting.PlayerA)},
		{CheckAction, actorToMove(acting.PlayerA), actorToMove(acting.PlayerB)},
		{CheckAction, actorToMove(acting.PlayerB), actorToMove(acting.ChanceId)},
		{DealPublicCardAction{&cards.QueenClubs}, actorToMove(acting.ChanceId), actorToMove(acting.PlayerA)},
		{BetAction, actorToMove(acting.PlayerA), actorToMove(acting.PlayerB)},
		{CallAction, actorToMove(acting.PlayerB), actorToMove(acting.ChanceId)},
		{DealPublicCardAction{&cards.JackHearts}, actorToMove(acting.ChanceId), actorToMove(acting.PlayerA)},
		{CheckAction, actorToMove(acting.PlayerA), actorToMove(acting.PlayerB)},
		{BetAction, actorToMove(acting.PlayerB), actorToMove(acting.PlayerA)},
		{CallAction, actorToMove(acting.PlayerA), noTest()},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfStacksChange(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs}, stackEqualsTo(acting.PlayerA, 100.), stackEqualsTo(acting.PlayerA, 100.-Ante)},
		{CheckAction, stackEqualsTo(acting.PlayerA, 100.-Ante), stackEqualsTo(acting.PlayerA, 100.-Ante)},
		{CheckAction, stackEqualsTo(acting.PlayerB, 100.-Ante), stackEqualsTo(acting.PlayerB, 100.-Ante)},
		{DealPublicCardAction{&cards.JackHearts}, stackEqualsTo(acting.PlayerA, 100.-Ante), stackEqualsTo(acting.PlayerA, 100.-Ante)},
		{BetAction, stackEqualsTo(acting.PlayerA, 100.-Ante), stackEqualsTo(acting.PlayerA, 100.-Ante-PostFlopBetSize)},
		{CallAction, stackEqualsTo(acting.PlayerB, 100.-Ante), stackEqualsTo(acting.PlayerB, 100.-Ante-PostFlopBetSize)},
		{DealPublicCardAction{&cards.QueenSpades}, stackEqualsTo(acting.PlayerB, 100.-Ante-PostFlopBetSize), stackEqualsTo(acting.PlayerB, 100.-Ante-PostFlopBetSize)},
		{CheckAction, stackEqualsTo(acting.PlayerA, 100.-Ante-PostFlopBetSize), stackEqualsTo(acting.PlayerA, 100.-Ante-PostFlopBetSize)},
		{BetAction, stackEqualsTo(acting.PlayerB, 100.-Ante-PostFlopBetSize), stackEqualsTo(acting.PlayerB, 100.-Ante-2*PostFlopBetSize)},
		{CallAction, stackEqualsTo(acting.PlayerA, 100.-Ante-PostFlopBetSize), stackEqualsTo(acting.PlayerA, 100.-Ante-2*PostFlopBetSize)},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfPotChanges(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs}, potEqualsTo(0.0), potEqualsTo(10.0)},
		{CheckAction, potEqualsTo(10), potEqualsTo(10)},
		{CheckAction, potEqualsTo(10), potEqualsTo(10)},
		{DealPublicCardAction{&cards.JackHearts}, potEqualsTo(10), potEqualsTo(10)},
		{BetAction, potEqualsTo(10), potEqualsTo(10 + PostFlopBetSize)},
		{CallAction, potEqualsTo(10 + PostFlopBetSize), potEqualsTo(10 + 2*PostFlopBetSize)},
		{DealPublicCardAction{&cards.QueenSpades}, potEqualsTo(10 + 2*PostFlopBetSize), potEqualsTo(10 + 2*PostFlopBetSize)},
		{CheckAction, potEqualsTo(10 + 2*PostFlopBetSize), potEqualsTo(10 + 2*PostFlopBetSize)},
		{BetAction, potEqualsTo(10 + 2*PostFlopBetSize), potEqualsTo(10 + 3*PostFlopBetSize)},
		{CallAction, potEqualsTo(10 + 3*PostFlopBetSize), potEqualsTo(10 + 4*PostFlopBetSize)},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestIfCardsGoToPlayers(t *testing.T) {

	root := createRootForTest(100., 100)

	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs}, noTest(), privateCards(cards.AceHearts, cards.KingClubs)},
		{CheckAction, noTest(), noTest()},
		{CheckAction, noTest(), noTest()},
		{DealPublicCardAction{&cards.JackHearts}, noTest(), flopCard(cards.JackHearts)},
		{CheckAction, noTest(), noTest()},
		{CheckAction, noTest(), noTest()},
		{DealPublicCardAction{&cards.QueenSpades}, noTest(), turnCard(cards.QueenSpades)},
		{CheckAction, noTest(), noTest()},
		{CheckAction, noTest(), noTest()},
	}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfChildPointersDifferFromParentsPointers(t *testing.T) {
	root := createRootForTest(100., 100.)
	child := root.Act(DealPrivateCardsAction{&cards.AceHearts, &cards.KingClubs})

	if child.(*RIGameState).actors[acting.ChanceId] == root.actors[acting.ChanceId] {
		t.Error("chance actor refers to the same value in both child and parent")
	}

	if child.(*RIGameState).actors[acting.PlayerA] == root.actors[acting.PlayerA] {
		t.Error("acting.PlayerA actor refers to the same value in both child and parent")
	}

	if child.(*RIGameState).actors[acting.PlayerB] == root.actors[acting.PlayerB] {
		t.Error("acting.PlayerB actor refers to the same value in both child and parent")
	}

	if child.(*RIGameState).table == root.table {
		t.Error("table should be different for child and parent")
	}
}

func TestGamePlayEvaluationFlushVsNothing(t *testing.T) {
	root := createRootForTest(100., 100.)
	hands := DealPrivateCardsAction{&cards.AceHearts, &cards.C2Spades}
	flop := DealPublicCardAction{&cards.JackHearts}
	turn := DealPublicCardAction{&cards.KingHearts}
	actions := []acting.Action{hands, CheckAction, CheckAction, flop, BetAction, CallAction, turn, CheckAction, BetAction, CallAction}

	singlePlayerPotContribution := Ante + PostFlopBetSize + PostFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationFlushVsStraightFlush(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.AceHearts, &cards.QueenHearts}
	flop := DealPublicCardAction{&cards.JackHearts}
	turn := DealPublicCardAction{&cards.KingHearts}

	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands, BetAction, CallAction, flop, BetAction, CallAction, turn, CheckAction, BetAction, CallAction}
	singlePlayerPotContribution := Ante + PreFlopBetSize + PostFlopBetSize + PostFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationPairVsPairDraw(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.AceHearts, &cards.AceSpades}
	flop := DealPublicCardAction{&cards.JackSpades}
	turn := DealPublicCardAction{&cards.KingHearts}

	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands, CheckAction, CheckAction, flop, CheckAction, CheckAction, turn, CheckAction, CheckAction}
	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(0), t)
}

func TestGamePlayEvaluationPairVsPairAWins(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.AceHearts, &cards.KingSpades}
	flop := DealPublicCardAction{&cards.AceSpades}
	turn := DealPublicCardAction{&cards.KingHearts}

	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands, CheckAction, CheckAction, flop, CheckAction, CheckAction, turn, CheckAction, CheckAction}
	singlePlayerPotContribution := Ante
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationPairVsPairBWinsBetterOwnCard(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.JackHearts, &cards.KingSpades}
	flop := DealPublicCardAction{&cards.C2Spades}
	turn := DealPublicCardAction{&cards.C2Hearts}

	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands, CheckAction, CheckAction, flop, CheckAction, CheckAction, turn, CheckAction, CheckAction}
	singlePlayerPotContribution := Ante
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationStraightVsStraightAWinsBetterOwnCard(t *testing.T) {

	hands := DealPrivateCardsAction{&cards.KingHearts, &cards.C10Spades}
	flop := DealPublicCardAction{&cards.JackClubs}
	turn := DealPublicCardAction{&cards.QueenDiamonds}

	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands, BetAction, CallAction, flop, CheckAction, CheckAction, turn, CheckAction, CheckAction}
	singlePlayerPotContribution := Ante + PreFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationPairVsThreeOfAKindBWins(t *testing.T) {

	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.KingClubs}
	flop := DealPublicCardAction{&cards.KingDiamonds}
	turn := DealPublicCardAction{&cards.KingSpades}

	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands, BetAction, CallAction, flop, CheckAction, CheckAction, turn, CheckAction, CheckAction}

	singlePlayerPotContribution := Ante + PreFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationOwnCardVsOwnCardAWins(t *testing.T) {

	hands := DealPrivateCardsAction{&cards.KingHearts, &cards.C10Clubs}
	flop := DealPublicCardAction{&cards.C2Diamonds}
	turn := DealPublicCardAction{&cards.C7Spades}

	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, BetAction, CallAction, flop, CheckAction, CheckAction, turn, BetAction, CallAction}

	singlePlayerPotContribution := Ante + PreFlopBetSize + PostFlopBetSize
	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAFoldsOnTurn(t *testing.T) {

	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.KingClubs}
	flop := DealPublicCardAction{&cards.KingDiamonds}
	turn := DealPublicCardAction{&cards.KingSpades}

	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, BetAction, CallAction, flop, CheckAction, CheckAction, turn, CheckAction, BetAction, FoldAction}

	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsOnTurn(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.KingClubs}
	flop := DealPublicCardAction{&cards.KingDiamonds}
	turn := DealPublicCardAction{&cards.KingSpades}

	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands, BetAction, CallAction, flop, CheckAction, CheckAction, turn, BetAction, FoldAction}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAFoldsOnFlop(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.KingClubs}
	flop := DealPublicCardAction{&cards.KingDiamonds}

	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, BetAction, CallAction, flop, CheckAction, BetAction, FoldAction}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsOnFlop(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.KingClubs}
	flop := DealPublicCardAction{&cards.KingDiamonds}

	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, BetAction, CallAction, flop, BetAction, FoldAction}
	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAFoldsPreFlop(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.KingClubs}
	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, CheckAction, BetAction, FoldAction}

	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsPreFlop(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.KingClubs}
	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, CheckAction, BetAction, RaiseAction, FoldAction}

	singlePlayerPotContribution := Ante + PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsPreFlopManyRaises(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.KingClubs}
	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, CheckAction, BetAction, RaiseAction, RaiseAction, RaiseAction, FoldAction}

	singlePlayerPotContribution := Ante + 3*PreFlopBetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)
}

func TestGamePlayInformationSetForAAfterRaises(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.KingClubs}
	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, CheckAction, BetAction, RaiseAction, RaiseAction}

	targetInformationSet := createInformationSet(cards.C10Hearts, cards.NoCard, cards.NoCard, actions)
	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForB_ChecksOnly(t *testing.T) {

	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.KingClubs}
	flop := DealPublicCardAction{&cards.KingDiamonds}
	turn := DealPublicCardAction{&cards.KingSpades}
	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands, CheckAction, CheckAction, flop, CheckAction, CheckAction, turn, CheckAction}
	targetInformationSet := createInformationSet(cards.KingClubs, cards.KingDiamonds, cards.KingSpades, actions)
	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForA_NoActions(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.KingClubs}

	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands}
	targetInformationSet := createInformationSet(cards.C10Hearts, cards.NoCard, cards.NoCard, actions)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForB_SingleCheck(t *testing.T) {

	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.AceSpades}
	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands, CheckAction}
	targetInformationSet := createInformationSet(cards.AceSpades, cards.NoCard, cards.NoCard, actions)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForB_BetCallAndChecksOnly(t *testing.T) {

	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.KingClubs}
	flop := DealPublicCardAction{&cards.KingDiamonds}
	turn := DealPublicCardAction{&cards.KingSpades}
	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands, CheckAction, BetAction, CallAction, flop, CheckAction, CheckAction, turn, CheckAction}
	targetInformationSet := createInformationSet(cards.KingClubs, cards.KingDiamonds, cards.KingSpades, actions)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForBAfterCheckBetRaise(t *testing.T) {

	hands := DealPrivateCardsAction{&cards.C10Hearts, &cards.QueenClubs}
	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands, CheckAction, BetAction, RaiseAction}
	targetInformationSet := createInformationSet(cards.QueenClubs, cards.NoCard, cards.NoCard, actions)

	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}
func testGamePlayAfterEveryAction(node *RIGameState, actionsTests []ActionTestsTriple, t *testing.T) {
	nodes := []games.GameState{node}
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

func testGamePlayAfterAllActions(node *RIGameState, actions []acting.Action, test func(state *RIGameState) bool, t *testing.T) {
	nodes := []games.GameState{node}
	for i := range actions {
		child := nodes[i].Act(actions[i])
		nodes = append(nodes, child)
	}
	if !test(nodes[len(nodes)-1].(*RIGameState)) {
		t.Error("post game test function did not pass")
	}
}

func createRootForTest(playerAStack float32, playerBStack float32) *RIGameState {
	playerA := &Player{Id: acting.PlayerA, Actions: nil, Card: nil, Stack: playerAStack}
	playerB := &Player{Id: acting.PlayerB, Actions: nil, Card: nil, Stack: playerBStack}
	return Root(playerA, playerB, cards.CreateFullDeck(true))
}

func roundCheck(expectedRound rounds.PokerRound) func(node *RIGameState) bool {
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
			if m.Name() == acting.Raise {
				return false
			}
		}
		return true
	}
}

func actorToMove(actorId acting.ActorID) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return state.nextToMove == actorId
	}
}

func stackEqualsTo(player acting.ActorID, stack float32) func(state *RIGameState) bool {
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
		if len(Actions) == 1 && Actions[0].Name() == acting.Check {
			return true
		}
		return false
	}
}
func checkAndBetAvailable() func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		Actions := state.Actions()
		if len(Actions) == 2 && Actions[0].Name() == acting.Check && Actions[1].Name() == acting.Bet {
			return true
		}
		return false
	}
}

func privateCards(playerACard cards.Card, playerBCard cards.Card) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return *(state.actors[acting.PlayerA].(*Player).Card) == playerACard && *(state.actors[acting.PlayerB].(*Player).Card) == playerBCard
	}
}

func flopCard(publicFlopCard cards.Card) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return state.table.Cards[0] == publicFlopCard
	}
}

func turnCard(publicTurnCard cards.Card) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		return state.table.Cards[1] == publicTurnCard
	}
}

func lastInformationSet(informationSet [InformationSetSizeBytes]byte) func(state *RIGameState) bool {
	return func(state *RIGameState) bool {
		currentInformationSet := state.InformationSet()
		return currentInformationSet == informationSet
	}
}

func createInformationSet(prvCard cards.Card, flopCard cards.Card, turnCard cards.Card, actions []acting.Action) [InformationSetSizeBytes]byte {

	informationSet := [InformationSetSizeBytes]byte{}
	informationSetBool := [InformationSetSize]bool{
		prvCard.Symbol[0], prvCard.Symbol[1], prvCard.Symbol[2], prvCard.Symbol[3],
		prvCard.Suit[0], prvCard.Suit[1], prvCard.Suit[2],
		flopCard.Symbol[0], flopCard.Symbol[1], flopCard.Symbol[2], flopCard.Symbol[3],
		flopCard.Suit[0], flopCard.Suit[1], flopCard.Suit[2],
		turnCard.Symbol[0], turnCard.Symbol[1], turnCard.Symbol[2], turnCard.Symbol[3],
		turnCard.Suit[0], turnCard.Suit[1], turnCard.Suit[2],
	}

	var currentAction acting.Action
	for i := 21; len(actions) > 0; i += 3 {
		// somehow tricky pop..
		currentAction, actions = actions[len(actions)-1], actions[:len(actions)-1]
		informationSetBool[i] = currentAction.Name()[0]
		informationSetBool[i+1] = currentAction.Name()[1]
		informationSetBool[i+2] = currentAction.Name()[2]
	}

	for i := 0; i < InformationSetSizeBytes; i++ {
		informationSet[i] = acting.CreateByte(informationSetBool[(i * 8):((i + 1) * 8)])
	}
	return informationSet
}
