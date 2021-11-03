package model

import (
	"fmt"
)

type locationArray []Location

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

//TODO - get json to serialize using strings below.
//TODO - look for a way that avoids repeating ourselves.
func (d Direction) String() string {
	return [...]string{"North", "South", "East", "West"}[d]
}

type Plan struct {
	Locations []Location
}

type LocationLabel string

type LocationId int
type Location struct {
	ID          LocationId
	Label       LocationLabel
	Transitions []Transition
}
func (location Location) String() string {
	directions := func(location Location) []string {
		var transitions []string
		for key, value := range location.Transitions {
			transitions = append(transitions, fmt.Sprintf("(%v) => %v", key, value))
		}
		return transitions
	}
	return fmt.Sprintf("%v: %q", location.Label, directions(location))
}

type TransitionId int
type Transition struct {
	ID 			TransitionId
	LocationId  LocationId
	Direction   Direction
	Destination LocationLabel
}

