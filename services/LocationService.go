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

func (service LocationService) Display() {
	tx := service.DB.MustBegin()
	service.locationDao.Map(*tx)
	locations := service.locationDao.GetAll(*tx)
	for _, location := range locations {
		fmt.Println(location.Label)
	}
	transitions := service.transitionDao.GetAll(*tx)
	for _, transitions := range transitions {
		fmt.Println(transitions)
	}
	
	tx.Commit()
	fmt.Println(locations)
}
