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
}
type Dao struct{
	Location LocationDao
	Transition TransitionDao
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
func (LocationDao) Get(tx sqlx.Tx) []Location {
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

func main() {
	//gormExample()
	//barr, _ := json.MarshalIndent(locations, "", "    ")
	//fmt.Println(string(barr))
	var plan Plan
	err := plan.LoadJson("3x3_floor")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sqlx.Connect("postgres", "user=test dbname=test sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	var dao Dao

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec("DROP OWNED BY test")
	db.MustExec(schema)
	tx := db.MustBegin()
	for i, location := range plan.Locations {
		fmt.Printf("%v %v\n", i, location)
		locationId := dao.Location.Insert(*tx, location)
		for _, transition := range location.Transitions {
			dao.Transition.Insert(*tx, locationId, transition)
		}
		fmt.Printf("%v", locationId)
	}
	dao.Location.Map(*tx)
	locations := dao.Location.Get(*tx)
	for _, location := range locations {
		fmt.Println(location.Label)
	}
	tx.Commit()

}

//var schema = `
//CREATE TABLE person (
//    first_name text,
//    last_name text,
//    email text
//);
//
//CREATE TABLE place (
//    country text,
//    city text NULL,
//    telcode integer
//)`
//
//type Person struct {
//	FirstName string `db:"first_name"`
//	LastName  string `db:"last_name"`
//	Email     string
//}
//
//type Place struct {
//	Country string
//	City    sql.NullString
//	TelCode int
//}
//
//func dbmain() {
//	// this Pings the database trying to connect
//	// use sqlx.Open() for sql.Open() semantics
//	db, err := sqlx.Connect("postgres", "user=test dbname=test sslmode=disable")
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	// exec the schema or fail; multi-statement Exec behavior varies between
//	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
//	db.MustExec(schema)
//
//	tx := db.MustBegin()
//	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)",
//		"Jason", "Moiron", "jmoiron@jmoiron.net")
//	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)",
//		"John", "Doe", "johndoeDNE@gmail.net")
//	tx.MustExec("INSERT INTO place (country, city, telcode) VALUES ($1, $2, $3)",
//		"United States", "New York", "1")
//	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)",
//		"Hong Kong", "852")
//	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)",
//		"Singapore", "65")
//	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
//	tx.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)",
//		&Person{"Jane", "Citizen", "jane.citzen@example.com"})
//	tx.Commit()
//
//	// Query the database, storing results in a []Person (wrapped in []interface{})
//	people := []Person{}
//	db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
//	jason, john := people[0], people[1]
//
//	fmt.Printf("%#v\n%#v", jason, john)
//	// Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}
//	// Person{FirstName:"John", LastName:"Doe", Email:"johndoeDNE@gmail.net"}
//
//	// You can also get a single result, a la QueryRow
//	jason = Person{}
//	err = db.Get(&jason, "SELECT * FROM person WHERE first_name=$1", "Jason")
//	fmt.Printf("%#v\n", jason)
//	// Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}
//
//	// if you have null fields and use SELECT *, you must use sql.Null* in your struct
//	places := []Place{}
//	err = db.Select(&places, "SELECT * FROM place ORDER BY telcode ASC")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	usa, singsing, honkers := places[0], places[1], places[2]
//
//	fmt.Printf("%#v\n%#v\n%#v\n", usa, singsing, honkers)
//	// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
//	// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}
//	// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}
//
//	// Loop through rows using only one struct
//	place := Place{}
//	rows, err := db.Queryx("SELECT * FROM place")
//	for rows.Next() {
//		err := rows.StructScan(&place)
//		if err != nil {
//			log.Fatalln(err)
//		}
//		fmt.Printf("%#v\n", place)
//	}
//	// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
//	// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}
//	// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}
//
//	// Named queries, using `:name` as the bindvar.  Automatic bindvar support
//	// which takes into account the dbtype based on the driverName on sqlx.Open/Connect
//	_, err = db.NamedExec(`INSERT INTO person (first_name,last_name,email) VALUES (:first,:last,:email)`,
//		map[string]interface{}{
//			"first": "Bin",
//			"last": "Smuth",
//			"email": "bensmith@allblacks.nz",
//		})
//
//	// Selects Mr. Smith from the database
//	rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:fn`, map[string]interface{}{"fn": "Bin"})
//
//	// Named queries can also use structs.  Their bind names follow the same rules
//	// as the name -> db mapping, so struct fields are lowercased and the `db` tag
//	// is taken into consideration.
//	rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:first_name`, jason)
//
//
//	// batch insert
//
//	// batch insert with structs
//	personStructs := []Person{
//		{FirstName: "Ardie", LastName: "Savea", Email: "asavea@ab.co.nz"},
//		{FirstName: "Sonny Bill", LastName: "Williams", Email: "sbw@ab.co.nz"},
//		{FirstName: "Ngani", LastName: "Laumape", Email: "nlaumape@ab.co.nz"},
//	}
//
//	_, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
//        VALUES (:first_name, :last_name, :email)`, personStructs)
//
//	// batch insert with maps
//	personMaps := []map[string]interface{}{
//		{"first_name": "Ardie", "last_name": "Savea", "email": "asavea@ab.co.nz"},
//		{"first_name": "Sonny Bill", "last_name": "Williams", "email": "sbw@ab.co.nz"},
//		{"first_name": "Ngani", "last_name": "Laumape", "email": "nlaumape@ab.co.nz"},
//	}
//
//	_, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
//        VALUES (:first_name, :last_name, :email)`, personMaps)
//}
///*
//package main
//
//func main() {
//g := &Game{Health: 100, Welcome: "Welcome to the Starship Enterprise\n\n", CurrentLocation: "Bridge"}
//g.Play()
//}
//*/
