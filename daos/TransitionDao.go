package daos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"takeoff.com/matilda/model"
)

type ITransitionDao interface {
	GetAll(tx sqlx.Tx) []model.Transition
	Insert(tx sqlx.Tx, locationId int, transition model.Transition)
}

type TransitionDao struct {}

func (TransitionDao) GetAll(tx sqlx.Tx) []model.Transition {
	var transitions []model.Transition
	err := tx.Select(&transitions, "SELECT * FROM transition ORDER BY locationId ASC")
	if err != nil {
		fmt.Printf("An error occured reading %s", err)
	}
	return transitions
}

func (TransitionDao) Insert(tx sqlx.Tx, locationId int, transition model.Transition) {
	tx.MustExec("INSERT INTO transition (locationId, direction, destination) VALUES ($1, $2, $3)",
		locationId, transition.Direction, transition.Destination)
}

