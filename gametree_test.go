package gocfr

import (
	"testing"
)

type MoveTestsTriple struct {
	move Move
	preTest func(state RhodeIslandGameState) bool
	postTest func(state RhodeIslandGameState) bool
}

func TestGameCreation(t *testing.T) {
	deck := CreateFullDeck()
	root := RhodeIslandGameState{Start, &deck, nil, nil, nil}
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

func TestGamePlay_1(t *testing.T) {
	deck := CreateFullDeck()
	root := RhodeIslandGameState{Start, &deck, nil, nil, nil}

	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, roundCheckFunc(Start), roundCheckFunc(PreFlop)},
		{Check,roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Check, roundCheckFunc(PreFlop),roundCheckFunc(PreFlop)},
		{DealPublicCard, roundCheckFunc(PreFlop), roundCheckFunc(Flop)},
		{Bet, roundCheckFunc(Flop),roundCheckFunc(Flop)},
		{Call, roundCheckFunc(Flop),roundCheckFunc(Flop)},
		{DealPublicCard,roundCheckFunc(Flop), roundCheckFunc(Turn)},
		{Check, roundCheckFunc(Turn),roundCheckFunc(Turn)},
		{Bet, roundCheckFunc(Turn), roundCheckFunc(Turn)},
		{Call, roundCheckFunc(Turn),  GameEndFunc()},
	}

	testGamePlay(root, movesTestsPairs, t)
}

func TestGamePlay_Max6Raises(t *testing.T) {
	deck := CreateFullDeck()
	root := RhodeIslandGameState{Start, &deck, nil, nil, nil}

	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, roundCheckFunc(Start), roundCheckFunc(PreFlop)},
		{Bet,roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop),roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop),roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop),roundCheckFunc(PreFlop)},
		{Raise,roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop),roundCheckFunc(PreFlop)},
		{Fold, roundCheckFunc(PreFlop), GameEndFunc()},
	}

	testGamePlay(root, movesTestsPairs, t)

	movesTestsPairs = []MoveTestsTriple{
		{DealPrivateCards, roundCheckFunc(Start), roundCheckFunc(PreFlop)},
		{Bet,roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop),roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop),roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop),roundCheckFunc(PreFlop)},
		{Raise,roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop),NoRaiseAvailable()},
		{Fold, roundCheckFunc(PreFlop), GameEndFunc()},
	}

	testGamePlay(root, movesTestsPairs, t)


	movesTestsPairs = []MoveTestsTriple{
		{DealPrivateCards, roundCheckFunc(Start), roundCheckFunc(PreFlop)},
		{Bet,roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop),roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop),roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop),roundCheckFunc(PreFlop)},
		{Raise,roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop),NoRaiseAvailable()},
		{Call, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{DealPublicCard, roundCheckFunc(PreFlop), roundCheckFunc(Flop)},
		{Bet,roundCheckFunc(Flop), roundCheckFunc(Flop)},
		{Raise, roundCheckFunc(Flop),roundCheckFunc(Flop)},
		{Raise, roundCheckFunc(Flop), roundCheckFunc(Flop)},
		{Call, roundCheckFunc(Flop),roundCheckFunc(Flop)},
		{DealPublicCard, roundCheckFunc(Flop), roundCheckFunc(Turn)},
		{Check, roundCheckFunc(Turn), roundCheckFunc(Turn)},
		{Check, roundCheckFunc(Turn), GameEndFunc()},
	}

	testGamePlay(root, movesTestsPairs, t)

}


func TestGamePlay_CheckIfPlayerToMoveCorrect(t *testing.T) {
	deck := CreateFullDeck()
	root := RhodeIslandGameState{Start, &deck, nil, nil, nil}

	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, playerToMoveFunc(Environment), playerToMoveFunc(PlayerA)},
		{Check,playerToMoveFunc(PlayerA), playerToMoveFunc(PlayerB)},
		{Check, playerToMoveFunc(PlayerB),playerToMoveFunc(Environment)},
		{DealPublicCard, playerToMoveFunc(Environment), playerToMoveFunc(PlayerA)},
		{Bet, playerToMoveFunc(PlayerA),playerToMoveFunc(PlayerB)},
		{Call, playerToMoveFunc(PlayerB),playerToMoveFunc(Environment)},
		{DealPublicCard,playerToMoveFunc(Environment), playerToMoveFunc(PlayerA)},
		{Check, playerToMoveFunc(PlayerA),playerToMoveFunc(PlayerB)},
		{Bet, playerToMoveFunc(PlayerB), playerToMoveFunc(PlayerA)},
		{Call, playerToMoveFunc(PlayerA),  playerToMoveFunc(Environment)},
	}

	testGamePlay(root, movesTestsPairs, t)
}



func testGamePlay(node RhodeIslandGameState, movesTests []MoveTestsTriple, t *testing.T) {
	nodes := []RhodeIslandGameState{node}
	for i, _ := range movesTests {
		actions := nodes[i].GetAvailableActions()
		actionIndex := selectActionByMove(actions, movesTests[i].move)
		child := nodes[i].Play(actions[actionIndex], nil)
		nodes = append(nodes, child)

		if ! movesTests[i].preTest(nodes[i]) {
			t.Errorf("pre test function  #%v did not pass", i)
		}

		if ! movesTests[i].postTest(child) {
			t.Errorf("post test function  #%v did not pass", i)
		}
	}
}
