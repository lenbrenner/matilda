package services

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"takeoff.com/matilda/daos"
	"takeoff.com/matilda/model"
)

type LocationService struct {
	DB            *sqlx.DB           `inject:"Db"`
	locationDao   daos.LocationDao   `inject:LocationDao`
	transitionDao daos.TransitionDao `inject:LocationDao`
}

func (service LocationService) LoadFromFile(filename string) {
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	var plan model.Plan
	err := plan.LoadJson(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	service.LoadAll(plan.Locations)
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
	tx.Commit()
}
