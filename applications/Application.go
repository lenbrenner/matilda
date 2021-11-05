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
	binder := axon.NewBinder(
		axon.NewPackage(
			axon.Bind("Db").
				To().Factory(databaseFactory).
				WithArgs(axon.Args{os.Getenv("DB_INSTANCE_NAME")}),
			axon.Bind("Application").
				To().StructPtr(new(Application)),
			axon.Bind("LocationDao").
				To().StructPtr(new(daos.LocationDao)),
			axon.Bind("TransitionDao").
				To().StructPtr(new(daos.TransitionDao)),
			axon.Bind("LocationService").
				To().StructPtr(new(services.LocationService)),
			))
	injector := axon.NewInjector(binder)
	application = injector.GetStructPtr("Application").(*Application)
	application.LoadPlan("3x3_floor")
}

func Get() Application {
	return *application
}

func (app Application) LoadPlan(filename string) {
	plan, err := LoadJson(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app.LocationService.LoadAll(plan.Locations)
}

func LoadJson(filename string) (model.Plan, error) {
	var plan model.Plan
	dat, err := util.LoadPlan(fmt.Sprintf("maps/%s.json", filename))
	if err == nil {
		err := json.Unmarshal(dat, &plan)
		return plan, err
	} else {
		return plan, err
	}
}
