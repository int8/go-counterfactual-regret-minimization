package gopoker

type InformationSet [InformationSetSize]byte

type GameState interface {
	Parent() GameState
	Act(Action Action) GameState
	InformationSet() InformationSet
	Actions() []Action
	IsTerminal() bool
	CurrentActor() Actor
	Evaluate() float32
}
