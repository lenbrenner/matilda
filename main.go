/*
 */

//Create some locations
//Create some products
//Create some people
//Using locations build paths
//Using paths find intersections.
//This is how slotting is optimized: https://github.com/TakeoffTech/SDIA#sdia-solver-for-dynamic-inventory-allocation-former-continuous-slotting
package main

import (
	"fmt"
	"github.com/eddieowens/axon"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"takeoff.com/matilda/daos"
	"takeoff.com/matilda/model"
)

func pathsTo(location model.Location, destination model.Location) {
	//fmt.Printf("Searching for paths from %v to %v")
	//seen := map[LocationLabel]bool{}
	for direction, label := range location.Transitions {
		//node := locations[label]
		fmt.Println("%v, %v", direction, label)

	}
}

type DoSomethingService struct {
	DB *sqlx.DB `inject:"Db"`
}
func (service DoSomethingService) DoSomething() {
	var locationDao daos.LocationDao
	var transitionDao daos.TransitionDao

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
		locationId := locationDao.Insert(*tx, location)
		for _, transition := range location.Transitions {
			transitionDao.Insert(*tx, locationId, transition)
		}
		fmt.Printf("%v", locationId)
	}
	locationDao.Map(*tx)
	locations := locationDao.GetAll(*tx)
	
	for _, location := range locations {
		fmt.Println(location.Label)
	}
	tx.Commit()
}

//func xmain() {
//	//Todo - inject this through singleton manager
//	DB := DatabaseFactory()
//	Service := DoSomethingService{DB: *DB}
//	Service.DoSomething()
//}

func DatabaseFactory(_ axon.Injector, args axon.Args) axon.Instance {
	var schema = `
CREATE SEQUENCE location_id_seq;
CREATE TABLE location (
	id INT NOT NULL DEFAULT NEXTVAL('location_id_seq'), 
    label text
);

CREATE SEQUENCE transition_id_seq;
CREATE TABLE transition (
	id INT NOT NULL DEFAULT NEXTVAL('transition_id_seq'),
	location_id INT, 
    direction INT,
    destination text
);`

	db, err := sqlx.Connect("postgres", "user=test dbname=test sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	db.MustExec("DROP OWNED BY test")
	db.MustExec(schema)
	return axon.StructPtr(db)
}

func CarFactory(_ axon.Injector, args axon.Args) axon.Instance {
	fmt.Println("Hey, a new Car is being made!")
	return axon.StructPtr(
		&Car{
			LockCode: args.String(0),
		})
}

type Application struct {
	Car *Car `inject:"Car"` //Todo - Try for something like Car.InjectTag
	Service *DoSomethingService `inject:"DoSomethingService"`
}

//Todo - Read this https://tutorialedge.net/golang/the-go-init-function
func initApplication() *Application {
	binder := axon.NewBinder(axon.NewPackage(
		axon.Bind("Application").To().StructPtr(new(Application)),
		axon.Bind("Car").To().Factory(CarFactory).WithArgs(axon.Args{os.Getenv("CAR_LOCK_CODE")}),
		axon.Bind("Engine").To().StructPtr(new(Engine)),
		axon.Bind("FuelInjector").To().StructPtr(new(FuelInjector)),
		axon.Bind("Db").To().Factory(DatabaseFactory).WithoutArgs(),
		axon.Bind("DoSomethingService").To().StructPtr(new(DoSomethingService)),
	))
	injector := axon.NewInjector(binder)
	return injector.GetStructPtr("Application").(*Application)
}

func main() {
	var app = initApplication()
	app.Car.Start()
	fmt.Println(app.Car.LockCode)
	app.Service.DoSomething()
}

type Starter interface {
	Start()
}

type Car struct {
	LockCode string
	Engine Starter `inject:"Engine"`
}

func (c *Car) Start() {
	fmt.Println("Starting the Car!")
	c.Engine.Start()
}

type Engine struct {
	FuelInjector Starter `inject:"FuelInjector"`
}

func (e *Engine) Start() {
	fmt.Println("Starting the Engine!")
	e.FuelInjector.Start()
}

type FuelInjector struct {
}

func (*FuelInjector) Start() {
	fmt.Println("Starting the FuelInjector!")
}
