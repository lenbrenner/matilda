package floor

import (
	"fmt"
	"gorm.io/gorm"
)

type LocationLabel string

type Location struct {
	gorm.Model
	Label       LocationLabel
	Transitions map[Direction]LocationLabel
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
