package gopoker

type ActorId int8

type Actor interface {
	GetId() ActorId
}
