package games

import "github.com/int8/gopoker/acting"

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
