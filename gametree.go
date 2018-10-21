package gopoker

type InformationSet interface{}

type GameState interface {
	Parent() GameState
	Act(Action Action) GameState
	InformationSet() InformationSet
	Actions() []Action
	IsTerminal() bool
	CurrentActor() Actor
	Evaluate() float32
}
