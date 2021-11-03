package applications

import (
	"github.com/eddieowens/axon"
	_ "github.com/lib/pq"
	"os"
	"takeoff.com/matilda/daos"
	"takeoff.com/matilda/services"
)

type Application struct {
	Service *services.LocationService `inject:"LocationService"`
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
	return injector.GetStructPtr("Application").(*Application)
}
