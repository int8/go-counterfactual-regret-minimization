package cfr

import "math"
import . "github.com/int8/gopoker"

type StrategyMap map[InformationSet]map[ActionName]float64

type CfrComputingRoutine struct {
	sigmaSum   StrategyMap
	sigma      StrategyMap
	regretsSum StrategyMap
	root       GameState
}

func (routine *CfrComputingRoutine) cumulateCfrRegret(infSet InformationSet, action ActionName, value float64) {
	if _, ok := routine.regretsSum[infSet]; !ok {
		routine.regretsSum[infSet] = map[ActionName]float64{}
	}
	routine.regretsSum[infSet][action] += value
}

func (routine *CfrComputingRoutine) cumulateSigma(infSet InformationSet, action ActionName, value float64) {
	if _, ok := routine.sigmaSum[infSet]; !ok {
		routine.sigmaSum[infSet] = map[ActionName]float64{}
	}
	routine.sigmaSum[infSet][action] += value
}

func (routine *CfrComputingRoutine) ComputeNashEquilibriumViaCFR(iterations int) {
	for i := 0; i < iterations; i++ {
		routine.cfrUtilityRecursive(routine.root, 1, 1)
	}
}

func (routine *CfrComputingRoutine) updateSigma(infSet InformationSet) {
	if _, ok := routine.sigma[infSet]; !ok {
		routine.sigma[infSet] = map[ActionName]float64{}
	}

	regretSum := 0.
	for _, k := range routine.regretsSum[infSet] {
		regretSum += math.Max(k, 0.0)
	}
	for action := range routine.regretsSum[infSet] {
		if regretSum > 0.0 {
			routine.sigma[infSet][action] = math.Max(routine.regretsSum[infSet][action], 0.0) / regretSum
		} else {
			routine.sigma[infSet][action] = 1. / float64(len(routine.regretsSum[infSet]))
		}
	}
}

func (routine *CfrComputingRoutine) actionProbability(infSet InformationSet, action ActionName, nrOfActions int) float64 {
	if _, ok := routine.sigma[infSet]; !ok {
		return 1. / float64(nrOfActions)
	}
	return routine.sigma[infSet][action]
}

// this is still not ready - currently only computes chance sampling utility
//TODO: replace recursive approach with stack based approach - should run much faster
func (routine *CfrComputingRoutine) cfrUtilityRecursive(state GameState, reachA float64, reachB float64) float64 {

	childrenStateUtilities := map[ActionName]float64{}
	if state.IsTerminal() {
		return state.Evaluate()
	}

	if state.CurrentActor().GetId() == ChanceId {
		action := state.Actions()[0] // this is fine *practically* because our FullDeck is shuffled when created
		return routine.cfrUtilityRecursive(state.Act(action), reachA, reachB)
	}

	infSet := state.InformationSet()

	value := 0.0
	actions := state.Actions()
	for _, action := range actions {
		childReachA := 1.0
		childReachB := 1.0
		if state.CurrentActor().GetId() == PlayerA {
			childReachA = childReachA * routine.actionProbability(infSet, action.Name(), len(actions))
		} else {
			childReachB = childReachB * routine.actionProbability(infSet, action.Name(), len(actions))
		}

		childStateUtility := routine.cfrUtilityRecursive(state.Act(action), childReachA, childReachB)
		value += routine.sigma[infSet][action.Name()] * childStateUtility

		childrenStateUtilities[action.Name()] = childStateUtility
	}
	var cfrReach, reach float64
	if state.CurrentActor().GetId() == PlayerA {
		cfrReach, reach = reachB, reachA
	} else {
		cfrReach, reach = reachA, reachB
	}

	for _, action := range actions {
		actionCfrRegret := float64(state.CurrentActor().GetId()) * cfrReach * (childrenStateUtilities[action.Name()] - value)
		routine.cumulateCfrRegret(infSet, action.Name(), actionCfrRegret)
		routine.cumulateSigma(infSet, action.Name(), reach*routine.actionProbability(infSet, action.Name(), len(actions)))
	}

	routine.updateSigma(infSet)
	return value
}
