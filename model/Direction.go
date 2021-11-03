package model

type Direction int

//TODO - get json to serialize using strings below.
//TODO - look for a way that avoids repeating ourselves.
const (
	North Direction = iota
	South
	East
	West
)

func (d Direction) String() string {
	return [...]string{"North", "South", "East", "West"}[d]
}

