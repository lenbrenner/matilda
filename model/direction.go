package model

//Run `go generate` and stringer will apper in direction_string.go.
//https://arjunmahishi.medium.com/golang-stringer-ad01db65e306
//To install `go get -u -a golang.org/x/tools/cmd/stringer`
//go:generate stringer -type=Direction
type Direction int

const (
	North Direction = iota
	South
	East
	West
)
