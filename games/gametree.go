package games

import "github.com/int8/go-counterfactual-regret-minimization/acting"

type InformationSet interface{}

type GameState interface {
	Parent() GameState
	Act(Action acting.Action) GameState
	InformationSet() InformationSet
	Actions() []acting.Action
	IsTerminal() bool
	CurrentActor() acting.Actor
	Evaluate() float32
}
