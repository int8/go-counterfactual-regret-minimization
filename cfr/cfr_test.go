package cfr

import (
	"encoding/gob"
	"github.com/int8/go-counterfactual-regret-minimization/acting"
	"github.com/int8/go-counterfactual-regret-minimization/cards"
	"github.com/int8/go-counterfactual-regret-minimization/games/kuhn"
	"github.com/int8/go-counterfactual-regret-minimization/games/rhodeisland"
	"os"
	"testing"
)

func TestKuhnPokerNashEquilibriumMatchesExpectedUtility(t *testing.T) {
	root := createRootForKuhnPokerTest(1000., 1000.)
	routine := CreateComputingRoutine(root)
	ne := routine.ComputeNashEquilibriumViaCFR(50000)
	utility := computeUtility(root, ne)
	if utility > -0.05 || utility < -0.06 {
		t.Error("Unless you are extremely unlucky, something is wrong with your CFR implementation")
	}
}

func TestRhodeISlandPokerNashEquilibrium(t *testing.T) {

	rhodeisland.MaxRaises = 0
	root := createRootForRhodeIslandPokerTest(1000., 1000.)
	routine := CreateComputingRoutine(root)
	routine.ComputeNashEquilibriumViaCFR(10000)
}

func createRootForKuhnPokerTest(playerAStack float32, playerBStack float32) *kuhn.KuhnGameState {
	playerA := &kuhn.Player{Id: acting.PlayerA, Actions: nil, Card: nil, Stack: playerAStack}
	playerB := &kuhn.Player{Id: acting.PlayerB, Actions: nil, Card: nil, Stack: playerBStack}
	return kuhn.Root(playerA, playerB)
}

func createRootForRhodeIslandPokerTest(playerAStack float32, playerBStack float32) *rhodeisland.RIGameState {
	playerA := &rhodeisland.Player{Id: acting.PlayerA, Actions: nil, Card: nil, Stack: playerAStack}
	playerB := &rhodeisland.Player{Id: acting.PlayerB, Actions: nil, Card: nil, Stack: playerBStack}
	return rhodeisland.Root(playerA, playerB, cards.CreateLimitedDeck(cards.C10, true))
}

func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func readNashEquilibriumFromGob(filePath string) *StrategyMap {
	ne := new(StrategyMap)
	readGob(filePath, ne)
	return ne
}
