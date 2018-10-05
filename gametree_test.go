package gocfr

import (
	"testing"
)

type MoveTestsTriple struct {
	move     Move
	preTest  func(state RhodeIslandGameState) bool
	postTest func(state RhodeIslandGameState) bool
}

func TestGameCreation(t *testing.T) {
	root := CreateRoot(100, 100)
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

	moves := root.actors[root.nextToMove].GetAvailableMoves(&root)

	if moves == nil {
		t.Error("Game root should have one action available, no actions available")
	}

	if len(moves) != 1 {
		t.Errorf("Game root should have one action available, %v actions available", len(moves))
	}

}

func TestGamePlay_1(t *testing.T) {
	root := CreateRoot(100, 100)
	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, roundCheckFunc(Start), roundCheckFunc(PreFlop)},
		{Check, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Check, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{DealPublicCard, roundCheckFunc(PreFlop), roundCheckFunc(Flop)},
		{Bet, roundCheckFunc(Flop), roundCheckFunc(Flop)},
		{Call, roundCheckFunc(Flop), roundCheckFunc(Flop)},
		{DealPublicCard, roundCheckFunc(Flop), roundCheckFunc(Turn)},
		{Check, roundCheckFunc(Turn), roundCheckFunc(Turn)},
		{Bet, roundCheckFunc(Turn), roundCheckFunc(Turn)},
		{Call, roundCheckFunc(Turn), GameEndFunc()},
	}

	testGamePlay(root, movesTestsPairs, t)
}

func TestGamePlay_Max6Raises(t *testing.T) {
	root := CreateRoot(100, 100)

	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, roundCheckFunc(Start), roundCheckFunc(PreFlop)},
		{Bet, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Fold, roundCheckFunc(PreFlop), GameEndFunc()},
	}

	testGamePlay(root, movesTestsPairs, t)

	movesTestsPairs = []MoveTestsTriple{
		{DealPrivateCards, roundCheckFunc(Start), roundCheckFunc(PreFlop)},
		{Bet, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), NoRaiseAvailable()},
		{Fold, roundCheckFunc(PreFlop), GameEndFunc()},
	}

	testGamePlay(root, movesTestsPairs, t)

	movesTestsPairs = []MoveTestsTriple{
		{DealPrivateCards, roundCheckFunc(Start), roundCheckFunc(PreFlop)},
		{Bet, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{Raise, roundCheckFunc(PreFlop), NoRaiseAvailable()},
		{Call, roundCheckFunc(PreFlop), roundCheckFunc(PreFlop)},
		{DealPublicCard, roundCheckFunc(PreFlop), roundCheckFunc(Flop)},
		{Bet, roundCheckFunc(Flop), roundCheckFunc(Flop)},
		{Raise, roundCheckFunc(Flop), roundCheckFunc(Flop)},
		{Raise, roundCheckFunc(Flop), roundCheckFunc(Flop)},
		{Call, roundCheckFunc(Flop), roundCheckFunc(Flop)},
		{DealPublicCard, roundCheckFunc(Flop), roundCheckFunc(Turn)},
		{Check, roundCheckFunc(Turn), roundCheckFunc(Turn)},
		{Check, roundCheckFunc(Turn), GameEndFunc()},
	}

	testGamePlay(root, movesTestsPairs, t)

}

func TestGamePlay_CheckIfPlayerToMoveCorrect(t *testing.T) {
	root := CreateRoot(100, 100)
	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, ActionMakerToMoveFunc(ChanceId), ActionMakerToMoveFunc(PlayerA)},
		{Check, ActionMakerToMoveFunc(PlayerA), ActionMakerToMoveFunc(PlayerB)},
		{Check, ActionMakerToMoveFunc(PlayerB), ActionMakerToMoveFunc(ChanceId)},
		{DealPublicCard, ActionMakerToMoveFunc(ChanceId), ActionMakerToMoveFunc(PlayerA)},
		{Bet, ActionMakerToMoveFunc(PlayerA), ActionMakerToMoveFunc(PlayerB)},
		{Call, ActionMakerToMoveFunc(PlayerB), ActionMakerToMoveFunc(ChanceId)},
		{DealPublicCard, ActionMakerToMoveFunc(ChanceId), ActionMakerToMoveFunc(PlayerA)},
		{Check, ActionMakerToMoveFunc(PlayerA), ActionMakerToMoveFunc(PlayerB)},
		{Bet, ActionMakerToMoveFunc(PlayerB), ActionMakerToMoveFunc(PlayerA)},
		{Call, ActionMakerToMoveFunc(PlayerA), ActionMakerToMoveFunc(NoActionMaker)},
	}

	testGamePlay(root, movesTestsPairs, t)
}

func testGamePlay(node RhodeIslandGameState, movesTests []MoveTestsTriple, t *testing.T) {
	nodes := []RhodeIslandGameState{node}
	for i, _ := range movesTests {

		child := nodes[i].CurrentActor().Act(&nodes[i], movesTests[i].move)
		nodes = append(nodes, child)

		if !movesTests[i].preTest(nodes[i]) {
			t.Errorf("pre test function  #%v did not pass", i)
		}

		if !movesTests[i].postTest(child) {
			t.Errorf("post test function  #%v did not pass", i)
		}
	}
}
