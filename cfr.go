package gocfr

const InformationSetSize int = 64

type InformationSet [InformationSetSize]byte

type StrategyMap map[InformationSet]map[Action]float64

type GameState interface {
	Child(Action Action) GameState
	Actions() []Action
	IsChance() bool
	IsTerminal() bool
	CurrentActor() Actor
	Evaluate() float64
	CurrentInformationSet() InformationSet
}

type CfrComputingRoutine struct {
	sigmaSum   StrategyMap
	sigma      StrategyMap
	regretsSum StrategyMap
	root       GameState
}

func (routine *CfrComputingRoutine) ComputeNashEquilibriumViaCFR(iterations int) {
	for i := 0; i < iterations; i++ {
		routine.cfrUtilityRecursive(routine.root, 1, 1)
	}
}

// this is still not ready - currently only computes chance sampling utility
func (routine *CfrComputingRoutine) cfrUtilityRecursive(state GameState, reachA float64, reachB float64) float64 {

	childrenStateUtilities := map[Action]float64{}
	if state.IsTerminal() {
		return state.Evaluate()
	}

	if state.IsChance() {
		// TODO: make sure this is random - what if someones deck is not shuffled ?
		action := state.Actions()[0] // this is fine practically because our FullDeck is shuffled when created
		return routine.cfrUtilityRecursive(state.Child(action), reachA, reachB)
	}

	value := 0.0
	for _, action := range state.Actions() {
		actionProbability := 1. / float64(len(state.Actions())) // routine.sigma[state.CurrentInformationSet()][action]
		childStateUtility := routine.cfrUtilityRecursive(state.Child(action), reachA, reachB)
		value += actionProbability * childStateUtility
		childrenStateUtilities[action] = childStateUtility
	}

	return value
}
