package cfr

import (
	"fmt"
	"github.com/int8/gopoker"
	"github.com/int8/gopoker/kuhn"
	"github.com/int8/gopoker/rhodeisland"
	"testing"
)

func TestKuhnPokerNashEquilibriumMatchesExpectedUtility(t *testing.T) {
	root := createRootForKuhnPokerTest(1000., 1000.)
	routine := CfrComputingRoutine{root: root, regretsSum: StrategyMap{}, sigma: StrategyMap{}, sigmaSum: StrategyMap{}}
	ne := routine.ComputeNashEquilibriumViaCFR(5000, true)
	utility := computeUtility(root, ne)
	if utility > -0.05 || utility < -0.06 {
		t.Error("Unless you are extremelly unlucky, something is wrong with your CFR implementation")
	}
}

func TestRhodeIslandNashEquilibriumComputation(t *testing.T) {
	root := createRootForRhodeIslandPokerTest(1000., 1000.)
	routine := CfrComputingRoutine{root: root, regretsSum: StrategyMap{}, sigma: StrategyMap{}, sigmaSum: StrategyMap{}}
	ne := routine.ComputeNashEquilibriumViaCFR(100, true)
	fmt.Println(ne)
}

func createRootForKuhnPokerTest(playerAStack float32, playerBStack float32) *kuhn.KuhnGameState {
	playerA := &kuhn.Player{Id: gopoker.PlayerA, Actions: nil, Card: nil, Stack: playerAStack}
	playerB := &kuhn.Player{Id: gopoker.PlayerB, Actions: nil, Card: nil, Stack: playerBStack}
	return kuhn.Root(playerA, playerB)
}

func createRootForRhodeIslandPokerTest(playerAStack float32, playerBStack float32) *rhodeisland.RIGameState {
	playerA := &rhodeisland.Player{Id: gopoker.PlayerA, Actions: nil, Card: nil, Stack: playerAStack}
	playerB := &rhodeisland.Player{Id: gopoker.PlayerB, Actions: nil, Card: nil, Stack: playerBStack}
	return rhodeisland.Root(playerA, playerB)
}
