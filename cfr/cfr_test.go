package cfr

import (
	"encoding/gob"
	"github.com/int8/gopoker"
	"github.com/int8/gopoker/cards"
	"github.com/int8/gopoker/games/kuhn"
	"github.com/int8/gopoker/games/rhodeisland"
	"os"
	"testing"
)

//func TestKuhnPokerNashEquilibriumMatchesExpectedUtility(t *testing.T) {
//	root := createRootForKuhnPokerTest(1000., 1000.)
//	routine := CfrComputingRoutine{root: root, regretsSum: StrategyMap{}, sigma: StrategyMap{}, sigmaSum: StrategyMap{}}
//	ne := routine.ComputeNashEquilibriumViaCFR(50000, true)
//	utility := computeUtility(root, ne)
//	if utility > -0.05 || utility < -0.06 {
//		t.Error("Unless you are extremelly unlucky, something is wrong with your CFR implementation")
//	}
//}

func TestRhodeISlandPokerNashEquilibrium(t *testing.T) {

	rhodeisland.MaxRaises = 0
	root := createRootForRhodeIslandPokerTest(1000., 1000.)
	routine := CfrComputingRoutine{root: root, regretsSum: StrategyMap{}, sigma: StrategyMap{}, sigmaSum: StrategyMap{}}
	routine.ComputeNashEquilibriumViaCFR(100000, true)
	//for infSet := range ne {
	//	fmt.Fprintf(os.Stdout, "%v %v \n", rhodeisland.PrettyPrintInformationSet(infSet), ne[infSet])
	//}
}

func createRootForKuhnPokerTest(playerAStack float32, playerBStack float32) *kuhn.KuhnGameState {
	playerA := &kuhn.Player{Id: gopoker.PlayerA, Actions: nil, Card: nil, Stack: playerAStack}
	playerB := &kuhn.Player{Id: gopoker.PlayerB, Actions: nil, Card: nil, Stack: playerBStack}
	return kuhn.Root(playerA, playerB)
}

func createRootForRhodeIslandPokerTest(playerAStack float32, playerBStack float32) *rhodeisland.RIGameState {
	playerA := &rhodeisland.Player{Id: gopoker.PlayerA, Actions: nil, Card: nil, Stack: playerAStack}
	playerB := &rhodeisland.Player{Id: gopoker.PlayerB, Actions: nil, Card: nil, Stack: playerBStack}
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
