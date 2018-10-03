package gocfr

import (
	"testing"
)

type MoveTestPair struct {
	move Move
	test func(state RhodeIslandGameState) bool
}

func TestGameCreation(t *testing.T) {
	deck := CreateFullDeck()
	root := RhodeIslandGameState{Start, &deck, nil, nil}
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

	actions := root.GetAvailableActions()

	if actions == nil {
		t.Error("Game root should have one action available, no actions available")
	}

	if len(actions) != 1 {
		t.Errorf("Game root should have one action available, %v actions available", len(actions))
	}
}

func TestGamePlay(t *testing.T) {
	deck := CreateFullDeck()
	root := RhodeIslandGameState{Start, &deck, nil, nil}

	movesTestsPairs := []MoveTestPair{
		{DealPrivateCards, roundCheckFunc(PreFlop)},
		{Check, roundCheckFunc(PreFlop)},
		{Check, roundCheckFunc(PreFlop)},
		{DealPublicCard, roundCheckFunc(Flop)},
		{Bet, roundCheckFunc(Flop)},
		{Call, roundCheckFunc(Flop)},
		{DealPublicCard, roundCheckFunc(Turn)},
		{Check, roundCheckFunc(Turn)},
		{Bet, roundCheckFunc(Turn)},
		{Call, roundCheckFunc(Turn)},
	}
	testGamePlayPostMove(root, movesTestsPairs, t)

	movesTestsPairs = []MoveTestPair{
		{DealPrivateCards, roundCheckFunc(Start)},
		{Check, roundCheckFunc(PreFlop)},
		{Check, roundCheckFunc(PreFlop)},
		{DealPublicCard, roundCheckFunc(PreFlop)},
		{Bet, roundCheckFunc(Flop)},
		{Call, roundCheckFunc(Flop)},
		{DealPublicCard, roundCheckFunc(Flop)},
		{Check, roundCheckFunc(Turn)},
		{Bet, roundCheckFunc(Turn)},
		{Call, roundCheckFunc(Turn)},
	}

	testGamePlayPreMove(root, movesTestsPairs, t)
}

func testGamePlay(node RhodeIslandGameState, movesTestsPairs []MoveTestPair, t *testing.T, mode string) {
	nodes := []RhodeIslandGameState{node}
	for i, _ := range movesTestsPairs {
		actions := nodes[i].GetAvailableActions()
		actionIndex := selectActionByMove(actions, movesTestsPairs[i].move)
		child := nodes[i].Play(actions[actionIndex])
		nodes = append(nodes, child)

		if (mode == "post" && movesTestsPairs[i].test(child)) || (mode == "pre" && movesTestsPairs[i].test(nodes[i])) {
			t.Errorf("function numer %v did not pass", i)
		}

	}
}

func testGamePlayPreMove(node RhodeIslandGameState, movesTestsPairs []MoveTestPair, t *testing.T) {
	testGamePlay(node, movesTestsPairs, t, "pre")
}

func testGamePlayPostMove(node RhodeIslandGameState, movesTestsPairs []MoveTestPair, t *testing.T) {
	testGamePlay(node, movesTestsPairs, t, "post")
}

func roundCheckFunc(expectedRound Round) func(node RhodeIslandGameState) bool {
	return func(node RhodeIslandGameState) bool { return node.round != expectedRound }
}
