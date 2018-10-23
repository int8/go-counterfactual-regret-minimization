package kuhn

import (
	. "github.com/int8/gopoker"
	. "github.com/int8/gopoker/cards"
	"math"
	"testing"
)

type ActionTestsTriple struct {
	Action   Action
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

	if root.round != Start {
		t.Error("Initial round of the game should be Start")
	}

	if root.IsTerminal() == true {
		t.Error("Game root should not be terminal")
	}

	actions := root.Actions()

	if len(actions) != 3*2 {
		t.Errorf("Game root should have %v actions available, %v actions available", 3*2, len(actions))
	}

}

func TestIfParentsCorrect(t *testing.T) {
	root := createRootForTest(100., 100.)
	child := root.Act(DealPrivateCardsAction{&KingHearts, &JackHearts})
	if child.Parent() != root {
		t.Error("Root child should have root as a parent")
	}
}

func TestIfStackLimitsAvailableActions(t *testing.T) {
	root5 := createRootForTest(1., 1.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&KingHearts, &QueenHearts}, noTest(), noTest()},
		{CheckAction, onlyCheckAvailable(), noTest()},
		{CheckAction, onlyCheckAvailable(), gameEnd()},
	}
	testGamePlayAfterEveryAction(root5, actionsTestsPairs, t)

	root15 := createRootForTest(2., 2.)
	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&KingHearts, &QueenHearts}, noTest(), noTest()},
		{CheckAction, checkAndBetAvailable(), noTest()},
		{CheckAction, checkAndBetAvailable(), gameEnd()},
	}

	testGamePlayAfterEveryAction(root15, actionsTestsPairs, t)

}

func TestGamePlayAssertRounds(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&KingHearts, &QueenHearts}, roundCheck(Start), roundCheck(PreFlop)},
		{CheckAction, roundCheck(PreFlop), roundCheck(PreFlop)},
		{CheckAction, roundCheck(PreFlop), roundCheck(PreFlop)},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfPlayerToMoveCorrect(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&KingHearts, &QueenHearts}, actorToMove(ChanceId), actorToMove(PlayerA)},
		{CheckAction, actorToMove(PlayerA), actorToMove(PlayerB)},
		{CheckAction, actorToMove(PlayerB), gameEnd()},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&KingHearts, &QueenHearts}, actorToMove(ChanceId), actorToMove(PlayerA)},
		{CheckAction, actorToMove(PlayerA), actorToMove(PlayerB)},
		{BetAction, actorToMove(PlayerB), actorToMove(PlayerA)},
		{CallAction, actorToMove(PlayerA), gameEnd()},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfStacksChange(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&JackHearts, &KingHearts}, stackEqualsTo(PlayerA, 100.), stackEqualsTo(PlayerA, 100.-Ante)},
		{CheckAction, stackEqualsTo(PlayerA, 100.-Ante), stackEqualsTo(PlayerA, 100.-Ante)},
		{CheckAction, stackEqualsTo(PlayerB, 100.-Ante), stackEqualsTo(PlayerB, 100.-Ante)},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&JackHearts, &KingHearts}, stackEqualsTo(PlayerA, 100.), stackEqualsTo(PlayerA, 100.-Ante)},
		{BetAction, stackEqualsTo(PlayerA, 100.-Ante), stackEqualsTo(PlayerA, 100.-Ante-BetSize)},
		{CallAction, stackEqualsTo(PlayerB, 100.-Ante), stackEqualsTo(PlayerB, 100.-Ante-BetSize)},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfPotChanges(t *testing.T) {
	root := createRootForTest(100., 100.)
	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&QueenHearts, &KingHearts}, potEqualsTo(0.0), potEqualsTo(2.0)},
		{BetAction, potEqualsTo(2), potEqualsTo(3)},
		{CallAction, potEqualsTo(3), potEqualsTo(4)},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&JackHearts, &KingHearts}, potEqualsTo(0.0), potEqualsTo(2.0)},
		{CheckAction, potEqualsTo(2), potEqualsTo(2)},
		{CheckAction, potEqualsTo(2), potEqualsTo(2)},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestIfCardsGoToPlayers(t *testing.T) {

	root := createRootForTest(100., 100)

	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&QueenHearts, &KingHearts}, noTest(), privateCards(QueenHearts, KingHearts)},
		{CheckAction, noTest(), noTest()},
		{CheckAction, noTest(), gameEnd()}}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfChildPointersDifferFromParentsPointers(t *testing.T) {
	root := createRootForTest(100., 100.)
	child := root.Act(DealPrivateCardsAction{&KingHearts, &JackHearts})

	if child.(*KuhnGameState).actors[ChanceId] == root.actors[ChanceId] {
		t.Error("chance actor refers to the same value in both child and parent")
	}

	if child.(*KuhnGameState).actors[PlayerA] == root.actors[PlayerA] {
		t.Error("PlayerA actor refers to the same value in both child and parent")
	}

	if child.(*KuhnGameState).actors[PlayerB] == root.actors[PlayerB] {
		t.Error("PlayerB actor refers to the same value in both child and parent")
	}

	if child.(*KuhnGameState).table == root.table {
		t.Error("table should be different for child and parent")
	}
}

func TestGamePlayEvaluationBWinsCheckBetFold(t *testing.T) {
	hands := DealPrivateCardsAction{&QueenHearts, &JackHearts}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, CheckAction, BetAction, FoldAction}

	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerB, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationAWinsBetFold(t *testing.T) {
	hands := DealPrivateCardsAction{&QueenHearts, &JackHearts}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, BetAction, FoldAction}

	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerA, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)

}

func TestGamePlayEvaluationAWinsBetCall(t *testing.T) {
	hands := DealPrivateCardsAction{&QueenHearts, &JackHearts}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, BetAction, CallAction}

	singlePlayerPotContribution := Ante + BetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerA, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)

}

func TestGamePlayEvaluationBWinsBetCall(t *testing.T) {
	hands := DealPrivateCardsAction{&QueenHearts, &KingHearts}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, BetAction, CallAction}

	singlePlayerPotContribution := Ante + BetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerB, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)

}

func TestGamePlayEvaluationAWinsCheckBetCall(t *testing.T) {
	hands := DealPrivateCardsAction{&QueenHearts, &JackHearts}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, CheckAction, BetAction, CallAction}

	singlePlayerPotContribution := Ante + BetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerA, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)

}

func TestGamePlayEvaluationBWinsCheckBetCall(t *testing.T) {
	hands := DealPrivateCardsAction{&QueenHearts, &KingHearts}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, CheckAction, BetAction, CallAction}

	singlePlayerPotContribution := Ante + BetSize

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerB, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)

}

func TestGamePlayEvaluationAWinsCheckCheck(t *testing.T) {
	hands := DealPrivateCardsAction{&QueenHearts, &JackHearts}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, CheckAction, CheckAction}

	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerA, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(singlePlayerPotContribution), t)

}

func TestGamePlayEvaluationBWinsCheckCheck(t *testing.T) {
	hands := DealPrivateCardsAction{&QueenHearts, &KingHearts}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, CheckAction, CheckAction}

	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerB, 100.+singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackAfterEvaluationEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-singlePlayerPotContribution), t)

}

func TestGamePlayInformationSetForA_NoActions(t *testing.T) {
	hands := DealPrivateCardsAction{&QueenHearts, &KingHearts}

	root := createRootForTest(100., 100.)

	actions := []Action{hands}
	targetInformationSet := createInformationSet(QueenHearts, actions)
	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForB_SingleCheck(t *testing.T) {

	hands := DealPrivateCardsAction{&QueenHearts, &KingHearts}
	root := createRootForTest(100., 100.)

	actions := []Action{hands, CheckAction}

	targetInformationSet := createInformationSet(KingHearts, actions)
	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func testGamePlayAfterEveryAction(node *KuhnGameState, actionsTests []ActionTestsTriple, t *testing.T) {
	nodes := []GameState{node}
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

func testGamePlayAfterAllActions(node *KuhnGameState, actions []Action, test func(state *KuhnGameState) bool, t *testing.T) {
	nodes := []GameState{node}
	for i := range actions {
		child := nodes[i].Act(actions[i])
		nodes = append(nodes, child)
	}
	if !test(nodes[len(nodes)-1].(*KuhnGameState)) {
		t.Error("post game test function did not pass")
	}
}

func createRootForTest(playerAStack float32, playerBStack float32) *KuhnGameState {
	playerA := &Player{Id: PlayerA, Actions: nil, Card: nil, Stack: playerAStack}
	playerB := &Player{Id: PlayerB, Actions: nil, Card: nil, Stack: playerBStack}
	return Root(playerA, playerB)
}

func roundCheck(expectedRound Round) func(node *KuhnGameState) bool {
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

func actorToMove(actorId ActorID) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		return state.nextToMove == actorId
	}
}

func stackEqualsTo(player ActorID, stack float32) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		return math.Abs(float64(state.actors[player].(*Player).Stack-stack)) < 1e-9
	}
}

func stackAfterEvaluationEqualsTo(player ActorID, stack float32) func(state *KuhnGameState) bool {
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
		if len(Actions) == 1 && Actions[0].Name() == Check {
			return true
		}
		return false
	}
}
func checkAndBetAvailable() func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		Actions := state.Actions()
		if len(Actions) == 2 && Actions[0].Name() == Check && Actions[1].Name() == Bet {
			return true
		}
		return false
	}
}

func privateCards(playerACard Card, playerBCard Card) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		return *(state.actors[PlayerA].(*Player).Card) == playerACard && *(state.actors[PlayerB].(*Player).Card) == playerBCard
	}
}

func lastInformationSet(informationSet [InformationSetSizeBytes]byte) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		currentInformationSet := state.InformationSet()
		return currentInformationSet == informationSet
	}
}

func createInformationSet(card Card, actions []Action) [InformationSetSizeBytes]byte {

	informationSetBool := [InformationSetSize]bool{
		card.Symbol[0], card.Symbol[1], card.Symbol[2], card.Symbol[3],
		card.Suit[0], card.Suit[1], card.Suit[2],
	}
	var currentAction Action
	for i := 7; len(actions) > 0; i += 3 {
		// somehow tricky pop..
		currentAction, actions = actions[len(actions)-1], actions[:len(actions)-1]
		informationSetBool[i] = currentAction.Name()[0]
		informationSetBool[i+1] = currentAction.Name()[1]
		informationSetBool[i+2] = currentAction.Name()[2]
	}
	informationSet := [InformationSetSizeBytes]byte{}

	for i := 0; i < InformationSetSizeBytes; i++ {
		informationSet[i] = CreateByte(informationSetBool[(i * 8):((i + 1) * 8)])
	}

	return informationSet
}
