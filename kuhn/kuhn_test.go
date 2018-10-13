package kuhn

import (
	. "github.com/int8/gopoker"
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
	child := root.Act(DealPrivateCardsAction{&AceHearts, &KingClubs})
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
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, noTest(), noTest()},
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
		{DealPrivateCardsAction{&AceHearts, &KingClubs}, potEqualsTo(0.0), potEqualsTo(2.0)},
		{BetAction, potEqualsTo(2), potEqualsTo(3)},
		{CallAction, potEqualsTo(3), potEqualsTo(4)},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)

	actionsTestsPairs = []ActionTestsTriple{
		{DealPrivateCardsAction{&JackHearts, &AceHearts}, potEqualsTo(0.0), potEqualsTo(2.0)},
		{CheckAction, potEqualsTo(2), potEqualsTo(2)},
		{CheckAction, potEqualsTo(2), potEqualsTo(2)},
	}
	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestIfCardsGoToPlayers(t *testing.T) {

	root := createRootForTest(100., 100)

	actionsTestsPairs := []ActionTestsTriple{
		{DealPrivateCardsAction{&QueenHearts, &AceHearts}, noTest(), privateCards(QueenHearts, AceHearts)},
		{CheckAction, noTest(), noTest()},
		{CheckAction, noTest(), gameEnd()}}

	testGamePlayAfterEveryAction(root, actionsTestsPairs, t)
}

func TestGamePlay_CheckIfChildPointersDifferFromParentsPointers(t *testing.T) {
	root := createRootForTest(100., 100.)
	child := root.Act(DealPrivateCardsAction{&AceHearts, &KingClubs})

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

func TestGamePlayEvaluationAFoldsPreFlop(t *testing.T) {
	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, CheckAction, BetAction, FoldAction}

	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(-2*singlePlayerPotContribution), t)
}

func TestGamePlayEvaluationBFoldsPreFlop(t *testing.T) {
	hands := DealPrivateCardsAction{&C10Hearts, &KingClubs}
	root := createRootForTest(100., 100.)
	actions := []Action{hands, BetAction, FoldAction}

	singlePlayerPotContribution := Ante

	testGamePlayAfterAllActions(root, actions, gameEnd(), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerA, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, stackEqualsTo(PlayerB, 100.-singlePlayerPotContribution), t)
	testGamePlayAfterAllActions(root, actions, gameResult(2*singlePlayerPotContribution), t)
}

func TestGamePlayInformationSetForA_NoActions(t *testing.T) {
	hands := DealPrivateCardsAction{&QueenHearts, &AceHearts}

	root := createRootForTest(100., 100.)

	actions := []Action{hands}
	targetInformationSet := [InformationSetSize]byte{byte(QueenHearts.Name), byte(QueenHearts.Suit)}
	targetInformationSet[2] = byte(DealPrivateCards)
	testGamePlayAfterAllActions(root, actions, lastInformationSet(targetInformationSet), t)
}

func TestGamePlayInformationSetForB_SingleCheck(t *testing.T) {

	hands := DealPrivateCardsAction{&QueenHearts, &AceHearts}
	root := createRootForTest(100., 100.)

	actions := []Action{hands, CheckAction}

	targetInformationSet := [InformationSetSize]byte{byte(AceHearts.Name), byte(AceHearts.Suit)}
	targetInformationSet[2] = byte(Check)
	targetInformationSet[3] = byte(DealPrivateCards)

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

func createRootForTest(playerAStack float64, playerBStack float64) *KuhnGameState {
	playerA := &Player{id: PlayerA, actions: nil, card: nil, stack: playerAStack}
	playerB := &Player{id: PlayerB, actions: nil, card: nil, stack: playerBStack}
	return root(playerA, playerB)
}

func roundCheck(expectedRound Round) func(node *KuhnGameState) bool {
	return func(node *KuhnGameState) bool { return node.round == expectedRound }
}

func gameEnd() func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool { return state.IsTerminal() }
}

func gameResult(result float64) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		evaluation := state.Evaluate()
		return evaluation == result
	}
}

func actorToMove(actorId ActorId) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		return state.nextToMove == actorId
	}
}

func stackEqualsTo(player ActorId, stack float64) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		return math.Abs(state.actors[player].(*Player).stack-stack) < 1e-9
	}
}

func potEqualsTo(pot float64) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		return math.Abs(state.table.Pot-pot) < 1e-9
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
		return *(state.actors[PlayerA].(*Player).card) == playerACard && *(state.actors[PlayerB].(*Player).card) == playerBCard
	}
}

func lastInformationSet(informationSet [InformationSetSize]byte) func(state *KuhnGameState) bool {
	return func(state *KuhnGameState) bool {
		currentInformationSet := state.InformationSet()
		return currentInformationSet == informationSet
	}
}
