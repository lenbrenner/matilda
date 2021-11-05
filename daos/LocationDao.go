package daos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"takeoff.com/matilda/model"
)

type ILocationDao interface {
	Insert(tx sqlx.Tx, location model.Location) int
	GetAll(tx sqlx.Tx) []model.Location
	Map(tx sqlx.Tx)
}

type LocationDao struct {}

func (LocationDao) Insert(tx sqlx.Tx, location model.Location) int {
	stmt, err := tx.PrepareNamed("INSERT INTO location (label, latitude, longitude) VALUES (:label, :latitude, :longitude) RETURNING id")
	var id int
	err = stmt.Get(&id, location)
	if err != nil {
		log.Fatalln(err)
	}
	return id
}

func (LocationDao) GetAll(tx sqlx.Tx) []model.Location {
	var locations []model.Location
	tx.Select(&locations, "SELECT * FROM location ORDER BY label ASC")
	return locations
}

func (LocationDao) Map(tx sqlx.Tx) {
	// Loop through rows using only one struct
	location := model.Location{}
	rows, _ := tx.Queryx("SELECT * FROM location")
	for rows.Next() {
		err := rows.StructScan(&location)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%#v\n", location)
	}
}
