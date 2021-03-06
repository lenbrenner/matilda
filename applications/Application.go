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
	LocationService services.ILocationService `inject:"LocationService"`
}

var application *Application

//This works in conjunction with `inject:"<Name>"` in our Daos and Services
//I can't use constants to keep bind and inject in sync:
//	https://github.com/golang/go/issues/40247
//Also I tried to use the new(<InstancePtr>) to derive name such as reflect.TypeOf(<InstancePtr>)
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
