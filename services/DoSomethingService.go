package services

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"takeoff.com/matilda/daos"
	"takeoff.com/matilda/model"
)

type DoSomethingService struct {
	DB            *sqlx.DB           `inject:"Db"`
	locationDao   daos.LocationDao   `inject:LocationDao`
	transitionDao daos.TransitionDao `inject:LocationDao`
}

func (service DoSomethingService) DoSomething() {
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	var plan model.Plan
	err := plan.LoadJson("3x3_floor")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tx := service.DB.MustBegin()
	for i, location := range plan.Locations {
		fmt.Printf("%v %v\n", i, location)
		locationId := service.locationDao.Insert(*tx, location)
		for _, transition := range location.Transitions {
			service.transitionDao.Insert(*tx, locationId, transition)
		}
		fmt.Printf("%v", locationId)
	}
	service.locationDao.Map(*tx)
	locations := service.locationDao.GetAll(*tx)

	for _, location := range locations {
		fmt.Println(location.Label)
	}
	tx.Commit()
}
