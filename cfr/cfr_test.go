package cfr

import (
	. "github.com/int8/gopoker"
	. "github.com/int8/gopoker/kuhn"
	"testing"
)

func TestKuhnPokerNashEquilibriumMatchesExpectedUtility(t *testing.T) {
	root := createRootForKuhnPokerTest(1000., 1000.)
	routine := CfrComputingRoutine{root: root, regretsSum: StrategyMap{}, sigma: StrategyMap{}, sigmaSum: StrategyMap{}}
	ne := routine.ComputeNashEquilibriumViaCFR(10000, true)
	utility := computeUtility(root, ne)
	if utility > -0.05 || utility < -0.06 {
		t.Error("Unless you are extremelly unlucky, something is wrong with your CFR implementation")
	}
}

func createRootForKuhnPokerTest(playerAStack float64, playerBStack float64) *KuhnGameState {
	playerA := &Player{Id: PlayerA, Actions: nil, Card: nil, Stack: playerAStack}
	playerB := &Player{Id: PlayerB, Actions: nil, Card: nil, Stack: playerBStack}
	return Root(playerA, playerB)
}
