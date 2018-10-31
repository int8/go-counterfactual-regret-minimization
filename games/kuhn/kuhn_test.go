package kuhn

import (
	"github.com/int8/gopoker/acting"
	"github.com/int8/gopoker/cards"
	"github.com/int8/gopoker/games"
	"github.com/int8/gopoker/rounds"
	"math"
	"testing"
)

type ActionTestsTriple struct {
	Action   acting.Action
	preTest  func(state *KuhnGameState) bool
	postTest func(state *KuhnGameState) bool
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

	if len(actions) != 3*2 {
		t.Errorf("Game root should have %v acting available, %v acting available", 3*2, len(actions))
	}

}

func TestIfParentsCorrect(t *testing.T) {
	root := createRootForTest(100., 100.)
	child := root.Act(DealPrivateCardsAction{&cards.KingHearts, &cards.JackHearts})
	if child.Parent() != root {
		t.Error("Root child should have root as a parent")
	}
}

func TestIfStackLimitsAvailableActions(t *testing.T) {
	root5 := createRootForTest(1., 1.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.KingHearts, &cards.QueenHearts}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
		{CheckAction, onlyCheckAvailable(), gameEnd()},
	}
	testGamePlayAfterEveryAction(root5, actionsTestsPairs, t)

	root15 := createRootForTest(2., 2.)
	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.KingHearts, &cards.QueenHearts}, noTest(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{CheckAction, checkAndBetAvailable(), gameEnd()},
	}

	testGamePlayAfterEveryAction(root15, actionsTestsPairs, t)

}

func TestGamePlayAssertRounds(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.KingHearts, &cards.QueenHearts}, roundCheck(rounds.Start), roundCheck(rounds.PreFlop)},
		{CheckAction, roundCheck(rounds.PreFlop), roundCheck(rounds.PreFlop)},
		{CheckAction, roundCheck(rounds.PreFlop), roundCheck(rounds.PreFlop)},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfPlayerToMoveCorrect(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.KingHearts, &cards.QueenHearts}, actorToMove(acting.ChanceId), actorToMove(acting.PlayerA)},
		{CheckAction, actorToMove(acting.PlayerA), actorToMove(acting.PlayerB)},
		{CheckAction, actorToMove(acting.PlayerB), gameEnd()},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.KingHearts, &cards.QueenHearts}, actorToMove(acting.ChanceId), actorToMove(acting.PlayerA)},
		{CheckAction, actorToMove(acting.PlayerA), actorToMove(acting.PlayerB)},
		{BetAction, actorToMove(acting.PlayerB), actorToMove(acting.PlayerA)},
		{CallAction, actorToMove(acting.PlayerA), gameEnd()},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfStacksChange(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.JackHearts, &cards.KingHearts}, stackEqualsTo(acting.PlayerA, 100.), stackEqualsTo(acting.PlayerA, 100.-Ante)},
		{CheckAction, stackEqualsTo(acting.PlayerA, 100.-Ante), stackEqualsTo(acting.PlayerA, 100.-Ante)},
		{CheckAction, stackEqualsTo(acting.PlayerB, 100.-Ante), stackEqualsTo(acting.PlayerB, 100.-Ante)},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.JackHearts, &cards.KingHearts}, stackEqualsTo(acting.PlayerA, 100.), stackEqualsTo(acting.PlayerA, 100.-Ante)},
		{BetAction, stackEqualsTo(acting.PlayerA, 100.-Ante), stackEqualsTo(acting.PlayerA, 100.-Ante-BetSize)},
		{CallAction, stackEqualsTo(acting.PlayerB, 100.-Ante), stackEqualsTo(acting.PlayerB, 100.-Ante-BetSize)},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfPotChanges(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.QueenHearts, &cards.KingHearts}, potEqualsTo(0.0), potEqualsTo(2.0)},
		{BetAction, potEqualsTo(2), potEqualsTo(3)},
		{CallAction, potEqualsTo(3), potEqualsTo(4)},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.JackHearts, &cards.KingHearts}, potEqualsTo(0.0), potEqualsTo(2.0)},
		{CheckAction, potEqualsTo(2), potEqualsTo(2)},
		{CheckAction, potEqualsTo(2), potEqualsTo(2)},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestIfCardsGoToPlayers(t *testing.T) {

	root := createRootForTest(100., 100)

	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&cards.QueenHearts, &cards.KingHearts}, noTest(), privateCards(cards.QueenHearts, cards.KingHearts)},
		{CheckAction, noTest(), noTest()},
		{CheckAction, noTest(), gameEnd()}}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfChildPointersDifferFromParentsPointers(t *testing.T) {
	root := createRootForTest(100., 100.)
	child := root.Act(DealPrivateCardsAction{&cards.KingHearts, &cards.JackHearts})

	if child.(*KuhnGameState).actors[acting.ChanceId] == root.actors[acting.ChanceId] {
		t.Error("chance actor refers to the same value in both child and parent")
	}

	if child.(*KuhnGameState).actors[acting.PlayerA] == root.actors[acting.PlayerA] {
		t.Error("acting.PlayerA actor refers to the same value in both child and parent")
	}

	if child.(*KuhnGameState).actors[acting.PlayerB] == root.actors[acting.PlayerB] {
		t.Error("acting.PlayerB actor refers to the same value in both child and parent")
	}

	if child.(*KuhnGameState).table == root.table {
		t.Error("table should be different for child and parent")
	}
}

func TestGamePlayEvaluationBWinsCheckBetFold(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.QueenHearts, &cards.JackHearts}
	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, CheckAction, BetAction, FoldAction}

	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerB, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAWinsBetFold(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.QueenHearts, &cards.JackHearts}
	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, BetAction, FoldAction}

	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerA, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)

}

func TestGamePlayEvaluationAWinsBetCall(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.QueenHearts, &cards.JackHearts}
	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, BetAction, CallAction}

	singlePlayerPotContribution := Ante + BetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerA, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)

}

func TestGamePlayEvaluationBWinsBetCall(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.QueenHearts, &cards.KingHearts}
	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, BetAction, CallAction}

	singlePlayerPotContribution := Ante + BetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerB, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)

}

func TestGamePlayEvaluationAWinsCheckBetCall(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.QueenHearts, &cards.JackHearts}
	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, CheckAction, BetAction, CallAction}

	singlePlayerPotContribution := Ante + BetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerA, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)

}

func TestGamePlayEvaluationBWinsCheckBetCall(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.QueenHearts, &cards.KingHearts}
	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, CheckAction, BetAction, CallAction}

	singlePlayerPotContribution := Ante + BetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerB, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)

}

func TestGamePlayEvaluationAWinsCheckCheck(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.QueenHearts, &cards.JackHearts}
	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, CheckAction, CheckAction}

	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerA, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)

}

func TestGamePlayEvaluationBWinsCheckCheck(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.QueenHearts, &cards.KingHearts}
	root := createRootForTest(100., 100.)
	actions := []acting.Action{hands, CheckAction, CheckAction}

	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerB, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(acting.PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)

}

func TestGamePlayInformationSetForA_NoActions(t *testing.T) {
	hands := DealPrivateCardsAction{&cards.QueenHearts, &cards.KingHearts}

	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands}
	targetInformationSet := createInformationSet(cards.QueenHearts, actions)
	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForB_SingleCheck(t *testing.T) {

	hands := DealPrivateCardsAction{&cards.QueenHearts, &cards.KingHearts}
	root := createRootForTest(100., 100.)

	actions := []acting.Action{hands, CheckAction}

	targetInformationSet := createInformationSet(cards.KingHearts, actions)
	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func testGamePlayAfterEveryAction(node *KuhnGameState, actionsTests []ActionTestsTriple, t *testing.T) {
	nodes := []games.GameState{node}
	for i := range actionsTests {

		if !actionsTests[i].preTest(nodes[i].(*KuhnGameState)) {
			t.Errorf("pre action test function  #%v did not pass", i)
		}

		child := nodes[i].Act(actionsTests[i].Action)
		nodes = append(nodes, child)

		if !actionsTests[i].postTest(child.(*KuhnGameState)) {
			t.Errorf("post action test function  #%v did not pass", i)
		}
	}
}

func testGamePlayAfterAllActions(node *KuhnGameState, actions []acting.Action, test func(state *KuhnGameState) bool, t *testing.T) {
	nodes := []games.GameState{node}
	for i := range actions {
		child := nodes[i].Act(actions[i])
		nodes = append(nodes, child)
	}
	if !test(nodes[len(nodes)-1].(*KuhnGameState)) {
		t.Error("post game test function did not pass")
	}
}

func createRootForTest(PlayerAStack float32, PlayerBStack float32) *KuhnGameState {
	PlayerA := &Player{Id: acting.PlayerA, Actions: nil, Card: nil, Stack: PlayerAStack}
	PlayerB := &Player{Id: acting.PlayerB, Actions: nil, Card: nil, Stack: PlayerBStack}
	return Root(PlayerA, PlayerB)
}

func roundCheck(expectedRound rounds.PokerRound) func(node *KuhnGameState) bool {
	return func(node *KuhnGameState) bool { return node.round == expectedRound }
}

func gameEnd() func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool { return state.IsTerminal() }
}

func gameResult(result float32) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		evaluation := state.Evaluate()
		return evaluation == result
	}
}

func actorToMove(actorId acting.ActorID) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		return state.nextToMove == actorId
	}
}

func stackEqualsTo(player acting.ActorID, stack float32) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		return math.Abs(float64(state.actors[player].(*Player).Stack-stack)) < 1e-9
	}
}

func stackAfterEvaluationEqualsTo(player acting.ActorID, stack float32) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		state.Evaluate()
		return math.Abs(float64(state.actors[player].(*Player).Stack-stack)) < 1e-9
	}
}

func potEqualsTo(pot float32) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		return math.Abs(float64(state.table.Pot-pot)) < 1e-9
	}
}

func noTest() func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		return true
	}
}

func onlyCheckAvailable() func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		Actions := state.Actions()
		if len(Actions) == 1 && Actions[0].Name() == acting.Check {
			return true
		}
		return false
	}
}
func checkAndBetAvailable() func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		Actions := state.Actions()
		if len(Actions) == 2 && Actions[0].Name() == acting.Check && Actions[1].Name() == acting.Bet {
			return true
		}
		return false
	}
}

func privateCards(PlayerACard cards.Card, PlayerBCard cards.Card) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		return *(state.actors[acting.PlayerA].(*Player).Card) == PlayerACard && *(state.actors[acting.PlayerB].(*Player).Card) == PlayerBCard
	}
}

func lastInformationSet(informationSet [InformationSetSizeBytes]byte) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		currentInformationSet := state.InformationSet()
		return currentInformationSet == informationSet
	}
}

func createInformationSet(card cards.Card, actions []acting.Action) [InformationSetSizeBytes]byte {

	informationSetBool := [InformationSetSize]bool{
		card.Symbol[0], card.Symbol[1], card.Symbol[2], card.Symbol[3],
		card.Suit[0], card.Suit[1], card.Suit[2],
	}
	var currentAction acting.Action
	for i := 7; len(actions) > 0; i += 3 {
		// somehow tricky pop..
		currentAction, actions = actions[len(actions)-1], actions[:len(actions)-1]
		informationSetBool[i] = currentAction.Name()[0]
		informationSetBool[i+1] = currentAction.Name()[1]
		informationSetBool[i+2] = currentAction.Name()[2]
	}
	informationSet := [InformationSetSizeBytes]byte{}

	for i := 0; i < InformationSetSizeBytes; i++ {
		informationSet[i] = acting.CreateByte(informationSetBool[(i * 8):((i + 1) * 8)])
	}

	return informationSet
}
