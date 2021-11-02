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
	"encoding/json"
	"fmt"
	"github.com/eddieowens/axon"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"takeoff.com/matilda/util"
)

func pathsTo(location Location, destination Location) {
	//fmt.Printf("Searching for paths from %v to %v")
	//seen := map[LocationLabel]bool{}
	for direction, label := range location.Transitions {
		//node := locations[label]
		fmt.Println("%v, %v", direction, label)

	}
}

type locationArray []Location

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

//TODO - get json to serialize using strings below.
//TODO - look for a way that avoids repeating ourselves.
func (d Direction) String() string {
	return [...]string{"North", "South", "East", "West"}[d]
}

type Plan struct {
	Locations []Location
}

type LocationLabel string

type LocationId int
type Location struct {
	ID          LocationId
	Label       LocationLabel
	Transitions []Transition
}

type TransitionId int
type Transition struct {
	ID 			TransitionId
	LocationId  LocationId
	Direction   Direction
	Destination LocationLabel
}

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

func (location Location) String() string {
	directions := func(location Location) []string {
		var transitions []string
		for key, value := range location.Transitions {
			transitions = append(transitions, fmt.Sprintf("(%v) => %v", key, value))
		}
		return transitions
	}
	return fmt.Sprintf("%v: %q", location.Label, directions(location))
}

func (plan *Plan) LoadJson(filename string) error {
	dat, err := util.Resource(fmt.Sprintf("maps/%s.json", filename))
	if err == nil {
		err := json.Unmarshal(dat, &plan)
		return err
	} else {
		return err
	}
}

type Crud interface {
	Insert(sqlx.Tx, Location)
	Get(sqlx.Tx) []Location
}

type LocationDao struct {}
func (LocationDao) Insert(tx sqlx.Tx, location Location) int {
	stmt, err := tx.PrepareNamed("INSERT INTO location (label) VALUES (:label) RETURNING id")
	var id int
	err = stmt.Get(&id, location)
	if err != nil {
		log.Fatalln(err)
	}
	return id
}
func (LocationDao) GetAll(tx sqlx.Tx) []Location {
	locations := []Location{}
	tx.Select(&locations, "SELECT * FROM location ORDER BY label ASC")
	return locations
}
func (LocationDao) Map(tx sqlx.Tx) {
	// Loop through rows using only one struct
	location := Location{}
	rows, _ := tx.Queryx("SELECT id, label FROM location")
	for rows.Next() {
		err := rows.StructScan(&location)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%#v\n", location)
	}
}

type TransitionDao struct {}
func (TransitionDao) Insert(tx sqlx.Tx, locationId int, transition Transition) {
	tx.MustExec("INSERT INTO transition (location_id, direction, destination) VALUES ($1, $2, $3)",
		locationId, transition.Direction, transition.Destination)
}

type DoSomethingService struct {
	db sqlx.DB `inject:"db"`
}
func (service DoSomethingService) DoSomething() {
	var locationDao LocationDao
	var transitionDao TransitionDao

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	var plan Plan
	err := plan.LoadJson("3x3_floor")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tx := service.db.MustBegin()
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

func xmain() {
	//gormExample()
	//barr, _ := json.MarshalIndent(locations, "", "    ")
	//fmt.Println(string(barr))
	//Todo - inject this through singleton manager
	db, err := sqlx.Connect("postgres", "user=test dbname=test sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	db.MustExec("DROP OWNED BY test")
	db.MustExec(schema)

	service := DoSomethingService{db: *db}
	service.DoSomething()
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

func CarFactory(_ axon.Injector, args axon.Args) axon.Instance {
	fmt.Println("Hey, a new Car is being made!")
	return axon.StructPtr(
			&Car{
				LockCode: args.String(0),
			})
}

func main() {
	binder := axon.NewBinder(axon.NewPackage(
		axon.Bind("Car").To().Factory(CarFactory).WithArgs(axon.Args{os.Getenv("CAR_LOCK_CODE")}),
		axon.Bind("Engine").To().StructPtr(new(Engine)),
		axon.Bind("FuelInjector").To().StructPtr(new(FuelInjector)),
	))
	injector := axon.NewInjector(binder)
	car := injector.GetStructPtr("Car").(*Car)
	car.Start()
	fmt.Println(car.LockCode)
}