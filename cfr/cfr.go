package cfr

import (
	"github.com/int8/go-counterfactual-regret-minimization/acting"
	"github.com/int8/go-counterfactual-regret-minimization/games"
	"math/rand"
)

type StrategyMap map[games.InformationSet]map[acting.ActionName]float32

type ComputingRoutine struct {
	sigmaSum   StrategyMap
	sigma      StrategyMap
	regretsSum StrategyMap
	root       games.GameState
}


func CreateComputingRoutine(root games.GameState) *ComputingRoutine {
	routine := ComputingRoutine{root: root, regretsSum: StrategyMap{}, sigma: StrategyMap{}, sigmaSum: StrategyMap{}}
	return &routine
}

func (routine *ComputingRoutine) cumulateCfrRegret(infSet games.InformationSet, action acting.ActionName, value float32) {
	if _, ok := routine.regretsSum[infSet]; !ok {
		routine.regretsSum[infSet] = map[acting.ActionName]float32{}
	}
	routine.regretsSum[infSet][action] += value
}

func (routine *ComputingRoutine) cumulateSigma(infSet games.InformationSet, action acting.ActionName, value float32) {
	if _, ok := routine.sigmaSum[infSet]; !ok {
		routine.sigmaSum[infSet] = map[acting.ActionName]float32{}
	}
	routine.sigmaSum[infSet][action] += value
}

func (routine *ComputingRoutine) ComputeNashEquilibriumViaCFR(iterations int) StrategyMap {

	for i := 0; i < iterations; i++ {
		routine.cfrUtilityRecursive(routine.root, 1, 1)
	}
	return routine.computeNashEquilibriumBasedOnStrategySum()
}

func (routine *ComputingRoutine) updateSigma(infSet games.InformationSet) {
	if _, ok := routine.sigma[infSet]; !ok {
		routine.sigma[infSet] = map[acting.ActionName]float32{}
	}

	regretSum := float32(0.)
	for _, k := range routine.regretsSum[infSet] {
		regretSum += maxFloat32(k, 0.0)
	}
	for action := range routine.regretsSum[infSet] {
		if regretSum > 0.0 {
			routine.sigma[infSet][action] = maxFloat32(routine.regretsSum[infSet][action], 0.0) / regretSum
		} else {
			routine.sigma[infSet][action] = 1. / float32(len(routine.regretsSum[infSet]))
		}
	}
}

func (routine *ComputingRoutine) actionProbability(infSet games.InformationSet, action acting.ActionName, nrOfActions int) float32 {
	if _, ok := routine.sigma[infSet]; !ok {
		return 1. / float32(nrOfActions)
	}
	return routine.sigma[infSet][action]
}

func (routine *ComputingRoutine) cfrUtilityRecursive(state games.GameState, reachA float32, reachB float32) float32 {

	childrenStateUtilities := map[acting.ActionName]float32{}
	if state.IsTerminal() {
		return state.Evaluate()
	}

	if state.CurrentActor().GetID() == acting.ChanceId {
		actions := state.Actions()
		action := actions[rand.Intn(len(actions))]
		return routine.cfrUtilityRecursive(state.Act(action), reachA, reachB)
	}

	infSet := state.InformationSet()
	value := float32(0.0)
	actions := state.Actions()
	for _, action := range actions {
		childReachA := reachA
		childReachB := reachB
		prob := routine.actionProbability(infSet, action.Name(), len(actions))

		if state.CurrentActor().GetID() == acting.PlayerA {
			childReachA *= prob
		} else {
			childReachB *= prob
		}

		childStateUtility := routine.cfrUtilityRecursive(state.Act(action), childReachA, childReachB)
		value += prob * childStateUtility

		childrenStateUtilities[action.Name()] = childStateUtility
	}

	var cfrReach, reach float32
	if state.CurrentActor().GetID() == acting.PlayerA {
		cfrReach, reach = reachB, reachA
	} else {
		cfrReach, reach = reachA, reachB
	}

	for _, action := range actions {
		if cfrReach != 0 {
			actionCfrRegret := float32(state.CurrentActor().GetID()) * cfrReach * (childrenStateUtilities[action.Name()] - value)
			routine.cumulateCfrRegret(infSet, action.Name(), actionCfrRegret)
		}
		if reach != 0 {
			routine.cumulateSigma(infSet, action.Name(), reach*routine.actionProbability(infSet, action.Name(), len(actions)))
		}
	}
	if reach != 0 {
		routine.updateSigma(infSet)
	}

	return value
}

func (routine *ComputingRoutine) computeNashEquilibriumBasedOnStrategySum() StrategyMap {
	nashEquilibrium := StrategyMap{}
	for infSet := range routine.sigmaSum {
		nashEquilibrium[infSet] = map[acting.ActionName]float32{}
		infSetSigmaSum := float32(0.0)
		for action := range routine.sigmaSum[infSet] {
			infSetSigmaSum += routine.sigmaSum[infSet][action]
		}

		for action := range routine.sigmaSum[infSet] {
			nashEquilibrium[infSet][action] = routine.sigmaSum[infSet][action] / infSetSigmaSum
		}
	}
	return nashEquilibrium
}

func computeUtility(state games.GameState, sigma StrategyMap) float32 {

	if state.IsTerminal() {
		return state.Evaluate()
	}

	if state.CurrentActor().GetID() == acting.ChanceId {
		actions := state.Actions()
		eval := float32(0.0)
		for _, action := range actions {
			eval += (1. / float32(len(actions))) * computeUtility(state.Act(action), sigma)
		}
		return eval
	}

	infSet := state.InformationSet()

	value := float32(0.0)
	actions := state.Actions()
	for _, action := range actions {
		value += sigma[infSet][action.Name()] * computeUtility(state.Act(action), sigma)
	}
	return value
}
