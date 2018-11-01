package cfr

import (
	"github.com/int8/go-counterfactual-regret-minimization/acting"
	"github.com/int8/go-counterfactual-regret-minimization/games"
	"math/rand"
	"sync"
)

type StrategyMap struct {
	Value map[games.InformationSet]map[acting.ActionName]float32
	mutex *sync.Mutex
}

func newStrategyMap() StrategyMap {
	return StrategyMap{Value: map[games.InformationSet]map[acting.ActionName]float32{}, mutex: &sync.Mutex{}}
}

func (sm StrategyMap) initIfZero(infSet games.InformationSet) {
	if _, ok := sm.Value[infSet]; !ok {
		sm.Value[infSet] = map[acting.ActionName]float32{}
	}
}

func (sm StrategyMap) hasValue(infSet games.InformationSet) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	_, ok := sm.Value[infSet]
	return ok
}

func (sm StrategyMap) setValue(infSet games.InformationSet, action acting.ActionName, value float32) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.initIfZero(infSet)
	sm.Value[infSet][action] = value
}

func (sm StrategyMap) getValue(infSet games.InformationSet, action acting.ActionName) float32 {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.initIfZero(infSet)
	return sm.Value[infSet][action]
}

func (sm StrategyMap) getKeys(infSet games.InformationSet) []acting.ActionName {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	actions := []acting.ActionName{}
	for action := range sm.Value[infSet] {
		actions = append(actions, action)
	}
	return actions
}

func (sm StrategyMap) nrOfInfSets() int {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	return len(sm.Value)
}


func (sm StrategyMap) sumValuesForInformationSet(infSet games.InformationSet) float32 {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	regretSum := float32(0.)
	for _, k := range sm.Value[infSet] {
		regretSum += maxFloat32(k, 0.0)
	}
	return regretSum
}

type ComputingRoutine struct {
	sigmaSum   StrategyMap
	sigma      StrategyMap
	regretsSum StrategyMap
	root       games.GameState
}

func CreateComputingRoutine(root games.GameState) *ComputingRoutine {
	routine := ComputingRoutine{root: root, regretsSum: newStrategyMap(), sigma: newStrategyMap(), sigmaSum: newStrategyMap()}
	return &routine
}

func (routine *ComputingRoutine) cumulateCfrRegret(infSet games.InformationSet, action acting.ActionName, value float32) {
	currentValue := routine.regretsSum.getValue(infSet, action)
	routine.regretsSum.setValue(infSet, action, currentValue + value)
}

func (routine *ComputingRoutine) cumulateSigma(infSet games.InformationSet, action acting.ActionName, value float32) {
	currentValue := routine.sigmaSum.getValue(infSet, action)
	routine.sigmaSum.setValue(infSet, action, currentValue + value)
}

func (routine *ComputingRoutine) ComputeNashEquilibriumViaCFR(iterations int, numThreads int) StrategyMap {

	for i := 0; i < iterations / numThreads; i++ {
		group := &sync.WaitGroup{}
		for j := 0; j < numThreads; j++ {
			group.Add(1)
			go func() {
				routine.cfrUtilityRecursive(routine.root, 1, 1)
				group.Done()
			}()
		}
		group.Wait()
	}
	return routine.computeNashEquilibriumBasedOnStrategySum()
}

func (routine *ComputingRoutine) updateSigma(infSet games.InformationSet) {

	regretSum := routine.regretsSum.sumValuesForInformationSet(infSet)
	actions := routine.regretsSum.getKeys(infSet)
	for _, action := range actions {
		if regretSum > 0.0 {
			routine.sigma.setValue(infSet, action, maxFloat32(routine.regretsSum.getValue(infSet, action), 0.0) / regretSum)
		} else {
			routine.sigma.setValue(infSet, action, 1. / float32(len(actions)))
		}
	}
}

func (routine *ComputingRoutine) actionProbability(infSet games.InformationSet, action acting.ActionName, nrOfActions int) float32 {
	if !routine.sigma.hasValue(infSet) {
		return 1. / float32(nrOfActions)
	}
	return routine.sigma.getValue(infSet, action)
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
		if cfrReach > 0 {
			actionCfrRegret := float32(state.CurrentActor().GetID()) * cfrReach * (childrenStateUtilities[action.Name()] - value)
			routine.cumulateCfrRegret(infSet, action.Name(), actionCfrRegret)
		}
		if reach > 0 {
			routine.cumulateSigma(infSet, action.Name(), reach*routine.actionProbability(infSet, action.Name(), len(actions)))
		}
	}

	if cfrReach > 0 {
		routine.updateSigma(infSet)
	}

	return value
}

func (routine *ComputingRoutine) computeNashEquilibriumBasedOnStrategySum() StrategyMap {
	nashEquilibrium := newStrategyMap()
	for infSet := range routine.sigmaSum.Value {
		nashEquilibrium.Value[infSet] = map[acting.ActionName]float32{}
		infSetSigmaSum := float32(0.0)
		for action := range routine.sigmaSum.Value[infSet] {
			infSetSigmaSum += routine.sigmaSum.Value[infSet][action]
		}

		for action := range routine.sigmaSum.Value[infSet] {
			nashEquilibrium.Value[infSet][action] = routine.sigmaSum.Value[infSet][action] / infSetSigmaSum
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
		value += sigma.Value[infSet][action.Name()] * computeUtility(state.Act(action), sigma)
	}

	return value
}
