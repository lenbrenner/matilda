package services

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"takeoff.com/matilda/daos"
	"takeoff.com/matilda/model"
)

type ILocationService interface {
	LoadAll(locations []model.Location)
	GetAll() []model.Location
}

//Todo - find a new home
//Facilitates mocking
type IDB interface {
	MustBegin() *sqlx.Tx
}

type LocationService struct {
	DB            IDB            `inject:"Db"`
	LocationDao   daos.ILocationDao   `inject:"LocationDao"`
	TransitionDao daos.ITransitionDao `inject:"TransitionDao"`
}

func (service LocationService) LoadAll(locations []model.Location) {
	for _, location := range locations {
		service.Load(location)
	}
}
func (service LocationService) Load(location model.Location) {
	tx := service.DB.MustBegin()
	locationId := service.LocationDao.Insert(*tx, location)
	for _, transition := range location.Transitions {
		service.TransitionDao.Insert(*tx, locationId, transition)
	}
	tx.Commit()
}

type BinLength struct {
	Bin    int
	Length int
}
type GroupedTransitions struct {
	BinLengths []BinLength
	model.Transition
}

func last(arr []int) int {
	return arr[len(arr)-1]
}

//https://medium.com/@geisonfgfg/functional-go-bc116f4c96a4
func groupByLocation(transitions []model.Transition) ([]int, []int, []int) {
	a := transitions[0 : len(transitions)-1]
	b := transitions[1:]
	starts := make([]int, 0)
	ends := make([]int, 0)
	counts := make([]int, 0)
	starts = append(starts, 0)
	for i, ai := range a {
		bi := b[i]
		if ai.LocationId != bi.LocationId {
			ends = append(ends, i+1)
			counts = append(counts, last(ends)-last(starts))
			starts = append(starts, i+1)
		}
	}
	ends = append(ends, len(transitions))
	counts = append(counts, last(ends)-last(starts))
	return starts, ends, counts
}

func boundaries(transitions []model.Transition) []int {
	a := transitions[0 : len(transitions)-1]
	b := transitions[1:]
	boundaries := make([]int, 0)
	for i, ai := range a {
		bi := b[i]
		if ai.LocationId != bi.LocationId {
			boundaries = append(boundaries, i+1)
		}
	}
	boundaries = append(boundaries, len(transitions))
	return boundaries
}

func (service LocationService) GetAll() []model.Location {
	tx := service.DB.MustBegin()
	service.LocationDao.Map(*tx)
	locations := service.LocationDao.GetAll(*tx)
	for _, location := range locations {
		fmt.Println(location.Label)
	}
	//transitions := service.TransitionDao.GetAll(*tx)
	if tx != nil {
		tx.Commit()
	}
	return locations
}

func (service LocationService) Display() {
	tx := service.DB.MustBegin()
	transitions := service.TransitionDao.GetAll(*tx)
	tx.Commit()
	//Todo - tidy this up
	boundaries := boundaries(transitions)
	fmt.Println(boundaries)
	starts, ends, counts := groupByLocation(transitions)
	fmt.Println(starts)
	fmt.Println(ends)
	fmt.Println(counts)
	for i, start := range starts {
		end := ends[i]
		fmt.Println(transitions[start:end])
	}
}
