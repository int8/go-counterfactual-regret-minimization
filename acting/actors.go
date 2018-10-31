package acting

const (
	PlayerA  ActorID = 1
	PlayerB          = -PlayerA
	ChanceId         = 0
)

type ActorID int8

type Actor interface {
	GetID() ActorID
}
