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

// TODO: add side effects, strategies updates to make it compute Nash Equilibrium
func Utility(node GameState, strategiesSum StrategyMap, cfrRegretsSum StrategyMap) float64 {
	if node.IsTerminal() {
		return node.Evaluate()
	}

	if node.IsChance() {
		action := node.Actions()[0]
		return Utility(node.Child(action), strategiesSum, cfrRegretsSum)
	}

	value := 0.0
	for _, action := range node.Actions() {
		x := Utility(node.Child(action), strategiesSum, cfrRegretsSum)
		value += 1. / float64(len(node.Actions())) * x
	}
	return value
}
