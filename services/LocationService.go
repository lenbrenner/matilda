package services

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"takeoff.com/matilda/daos"
	"takeoff.com/matilda/model"
)

type LocationService struct {
	DB            *sqlx.DB           `inject:"Db"`
	locationDao   daos.LocationDao   `inject:LocationDao`
	transitionDao daos.TransitionDao `inject:TransitionDao`
}

func (service LocationService) LoadAll(locations []model.Location) {
	for _, location := range locations {
		service.Load(location)
	}
}
func (service LocationService) Load(location model.Location) {
	tx := service.DB.MustBegin()
	locationId := service.locationDao.Insert(*tx, location)
	for _, transition := range location.Transitions {
		service.transitionDao.Insert(*tx, locationId, transition)
	}
	tx.Commit()
}

type BinLength struct {
	Bin int
	Length int
}
type GroupedTransitions struct {
	BinLengths []BinLength
	model.Transition
}

func last(arr []int) int {
	return arr[len(arr) - 1]
}
//https://medium.com/@geisonfgfg/functional-go-bc116f4c96a4
func groupByLocation(transitions []model.Transition) ([]int, []int, []int){
	a := transitions[0:len(transitions)-1]
	b := transitions[1:]
	starts := make([]int, 0)
	ends := make([]int, 0)
	counts := make([]int, 0)
	starts = append(starts, 0)
	for i, ai := range(a) {
		bi := b[i]
		if ai.LocationId != bi.LocationId {
			ends = append(ends, i + 1)
			counts = append(counts, last(ends) - last(starts))
			starts = append(starts, i + 1)
		}
	}
	ends = append(ends, len(transitions))
	counts = append(counts, last(ends) - last(starts))
	return starts, ends, counts
}

func boundaries(transitions []model.Transition) []int {
	a := transitions[0:len(transitions)-1]
	b := transitions[1:]
	boundaries := make([]int, 0)
	for i, ai := range(a) {
		bi := b[i]
		if ai.LocationId != bi.LocationId {
			boundaries = append(boundaries, i + 1)
		}
	}
	boundaries = append(boundaries, len(transitions))
	return boundaries
}

func (service LocationService) Display() {
	tx := service.DB.MustBegin()
	service.locationDao.Map(*tx)
	locations := service.locationDao.GetAll(*tx)
	for _, location := range locations {
		fmt.Println(location.Label)
	}
	transitions := service.transitionDao.GetAll(*tx)
	//Todo - tidy this up
	boundaries := boundaries(transitions)
	fmt.Println(boundaries)
	starts, ends, counts := groupByLocation(transitions)
	fmt.Println(starts)
	fmt.Println(ends)
	fmt.Println(counts)
	for i, start := range(starts) {
		end := ends[i]
		fmt.Println(transitions[start:end])
	}
	tx.Commit()
}
