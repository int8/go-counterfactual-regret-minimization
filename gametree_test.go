package gocfr

import (
	"testing"
)

type MoveTestsTriple struct {
	move     Move
	preTest  func(state *GameState) bool
	postTest func(state *GameState) bool
}

func TestGameCreation(t *testing.T) {
	root := CreateRoot(PlayerA, 100., 100.)
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
	root := CreateRoot(PlayerA, 100., 100.)
	child := root.actors[root.nextToMove].(*Chance).Act(root, DealPrivateCards)
	if child.parent != root {
		t.Error("Root child should have root as a parent")
	}
}

func TestIfStackLimitsAvailableActions(t *testing.T) {
	root5 := CreateRoot(PlayerA, 5., 5.)
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
	testGamePlay(root5, movesTestsPairs, t)

	root15 := CreateRoot(PlayerA, 15., 15.)
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

	testGamePlay(root15, movesTestsPairs, t)

	root1000_15 := CreateRoot(PlayerA, 1000., 15.)
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

	testGamePlay(root1000_15, movesTestsPairs, t)

	root1000_35 := CreateRoot(PlayerA, 1000., 35.)
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

	testGamePlay(root1000_35, movesTestsPairs, t)
}

func TestGamePlay_1(t *testing.T) {
	root := CreateRoot(PlayerA, 100., 100.)
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

	testGamePlay(root, movesTestsPairs, t)
}

func TestGamePlay_Max6Raises(t *testing.T) {
	root := CreateRoot(PlayerA, 100., 100.)

	movesTestsPairs := []MoveTestsTriple{
		{DealPrivateCards, roundCheck(Start), roundCheck(PreFlop)},
		{Bet, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Raise, roundCheck(PreFlop), roundCheck(PreFlop)},
		{Fold, roundCheck(PreFlop), gameEnd()},
	}

	testGamePlay(root, movesTestsPairs, t)

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

	testGamePlay(root, movesTestsPairs, t)

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

	testGamePlay(root, movesTestsPairs, t)

}

func TestGamePlay_CheckIfPlayerToMoveCorrect(t *testing.T) {
	root := CreateRoot(PlayerA, 100., 100.)
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

	testGamePlay(root, movesTestsPairs, t)
}

func TestGamePlay_CheckIfStacksChange(t *testing.T) {
	root := CreateRoot(PlayerA, 100., 100.)
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

	testGamePlay(root, movesTestsPairs, t)
}

func TestIfRootCreationWithDeckPreparedWorks(t *testing.T) {
	aceHearts := Card{Ace, Hearts}
	c2Spades := Card{C2, Spades}
	jackHearts := Card{Jack, Hearts}
	kingHearts := Card{King, Hearts}

	preparedDeck := prepareDeckForTest(aceHearts, c2Spades, jackHearts, kingHearts)
	root := createRootWithPreparedDeck(PlayerA, 100., 100., preparedDeck)

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

	testGamePlay(root, movesTestsPairs, t)

}

func TestGamePlay_CheckIfChildPointersDifferFromParentsPointers(t *testing.T) {
	root := CreateRoot(PlayerA, 100., 100.)
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

func testGamePlay(node *GameState, movesTests []MoveTestsTriple, t *testing.T) {
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
