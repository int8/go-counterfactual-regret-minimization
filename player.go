package gopoker

type ActorID int8

type Actor interface {
	GetID() ActorID
}
