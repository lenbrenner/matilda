package applications

import (
	"github.com/eddieowens/axon"
	_ "github.com/lib/pq"
	"os"
	"takeoff.com/matilda/services"
)

type Application struct {
	Service *services.DoSomethingService `inject:"DoSomethingService"`
}

//Todo - Read this https://tutorialedge.net/golang/the-go-init-function
func InitApplication() *Application {
	binder := axon.NewBinder(axon.NewPackage(
		axon.Bind("Application").To().StructPtr(new(Application)),
		axon.Bind("Db").To().Factory(databaseFactory).WithArgs(axon.Args{os.Getenv("DB_INSTANCE_NAME")}),
		axon.Bind("DoSomethingService").To().StructPtr(new(services.DoSomethingService)),
	))
	injector := axon.NewInjector(binder)
	return injector.GetStructPtr("Application").(*Application)
}
