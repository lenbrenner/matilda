package applications

import (
	"fmt"
	"github.com/eddieowens/axon"
	_ "github.com/lib/pq"
	"os"
	"takeoff.com/matilda/daos"
	"takeoff.com/matilda/model"
	"takeoff.com/matilda/services"
)

type Application struct {
	LocationService *services.LocationService `inject:"LocationService"`
}

//Todo - Read this https://tutorialedge.net/golang/the-go-init-function
func InitApplication() *Application {
	binder := axon.NewBinder(axon.NewPackage(
		axon.Bind("Application").To().StructPtr(new(Application)),
		axon.Bind("Db").To().Factory(databaseFactory).WithArgs(axon.Args{os.Getenv("DB_INSTANCE_NAME")}),
		axon.Bind("LocationService").To().StructPtr(new(services.LocationService)),
		axon.Bind("LocationDao").To().StructPtr(new(daos.LocationDao)),
		axon.Bind("TransitionDao").To().StructPtr(new(daos.TransitionDao)),
	))
	injector := axon.NewInjector(binder)
	app := injector.GetStructPtr("Application").(*Application)
	app.LocationService.LoadFromFile("3x3_floor")
	return app
}

func (app Application) LoadPlan(filename string) {
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	var plan model.Plan
	err := plan.LoadJson(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app.LocationService.LoadAll(plan.Locations)
}
