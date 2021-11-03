package model

import "fmt"

type LocationId int

type LocationLabel string

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

