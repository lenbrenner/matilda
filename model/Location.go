package model

type LocationId int

type LocationLabel string

type Location struct {
	ID          LocationId
	Latitude    int32
	Longitude   int32
	Label       LocationLabel
	Transitions []Transition
}
