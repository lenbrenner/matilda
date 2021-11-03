package applications

import (
	"encoding/json"
	"fmt"
	"github.com/eddieowens/axon"
	_ "github.com/lib/pq"
	"os"
	"takeoff.com/matilda/daos"
	"takeoff.com/matilda/model"
	"takeoff.com/matilda/services"
	"takeoff.com/matilda/util"
)

type Application struct {
	LocationService *services.LocationService `inject:"LocationService"`
}

var application *Application
func init() {
	binder := axon.NewBinder(axon.NewPackage(
		axon.Bind("Application").
			To().StructPtr(new(Application)),
		axon.Bind("Db").
			To().Factory(databaseFactory).
			WithArgs(axon.Args{os.Getenv("DB_INSTANCE_NAME")}),
		axon.Bind("LocationService").
			To().StructPtr(new(services.LocationService)),
		axon.Bind("LocationDao").
			To().StructPtr(new(daos.LocationDao)),
		axon.Bind("TransitionDao").
			To().StructPtr(new(daos.TransitionDao)),
	))
	injector := axon.NewInjector(binder)
	application = injector.GetStructPtr("Application").(*Application)
	application.LoadPlan("3x3_floor")
}

func Get() Application {
	return *application
}

func (app Application) LoadPlan(filename string) {
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	var plan model.Plan
	err := LoadJson(filename, &plan)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app.LocationService.LoadAll(plan.Locations)
}

func LoadJson(filename string, plan *model.Plan) error {
	dat, err := util.LoadPlan(fmt.Sprintf("maps/%s.json", filename))
	if err == nil {
		err := json.Unmarshal(dat, &plan)
		return err
	} else {
		return err
	}
}
