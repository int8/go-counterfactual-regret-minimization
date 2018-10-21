package gopoker

type InformationSet [36 + 21]bool

type GameState interface {
	Parent() GameState
	Act(Action Action) GameState
	InformationSet() InformationSet
	Actions() []Action
	IsTerminal() bool
	CurrentActor() Actor
	Evaluate() float32
}

func DistanceToRoot(state GameState) int {
	if state.Parent() == nil {
		return 0
	}
	return 1 + DistanceToRoot(state.Parent())
}
