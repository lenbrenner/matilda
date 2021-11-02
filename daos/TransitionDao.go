package daos

import (
	"github.com/jmoiron/sqlx"
	"takeoff.com/matilda/model"
)

type TransitionDao struct {}
func (TransitionDao) Insert(tx sqlx.Tx, locationId int, transition model.Transition) {
	tx.MustExec("INSERT INTO transition (location_id, direction, destination) VALUES ($1, $2, $3)",
		locationId, transition.Direction, transition.Destination)
}

