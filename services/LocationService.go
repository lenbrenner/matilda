package services

import (
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
	DB            IDB                 `inject:"Db"`
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

func (service LocationService) GetAll() []model.Location {
	tx := service.DB.MustBegin()
	locations := service.LocationDao.GetAll(*tx)
	transitions := service.TransitionDao.GetAll(*tx)
	starts, ends, _ := encodeTransitions(transitions)
	transitionsByLocation := make(map[model.LocationId][]model.Transition, len(starts))
	//Todo - tidy this up
	for i, start := range starts {
		end := ends[i]
		transitionsByLocation[transitions[start].LocationId] = transitions[start:end]
	}
	for i, location := range locations {
		locations[i].Transitions = transitionsByLocation[location.ID]
	}
	if tx != nil {
		tx.Commit()
	}
	return locations
}

func last(arr []int) int {
	return arr[len(arr)-1]
}

func encodeTransitions(
	transitions []model.Transition,
	) ([]int, []int, []int) {
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
