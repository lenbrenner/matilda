package applications

import (
	"fmt"
	"github.com/eddieowens/axon"
	"github.com/jmoiron/sqlx"
	"log"
)

func databaseFactory(_ axon.Injector, args axon.Args) axon.Instance {
	var instanceName = args.String(0)
	fmt.Printf("Starting with database instance: %v\n", instanceName)
	var schema = `
CREATE SEQUENCE location_id_seq;
CREATE TABLE location (
	id INT NOT NULL DEFAULT NEXTVAL('location_id_seq'), 
    label text
);

CREATE SEQUENCE transition_id_seq;
CREATE TABLE transition (
	id INT NOT NULL DEFAULT NEXTVAL('transition_id_seq'),
	locationId INT, 
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
