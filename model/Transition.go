package model

type TransitionId int
type Transition struct {
	ID          TransitionId
	LocationId  LocationId
	Direction   Direction
	Destination LocationLabel
}
