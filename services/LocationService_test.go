package services

import (
	"fmt"
	"github.com/eddieowens/axon"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/mock"
	"takeoff.com/matilda/model"
	"testing"
)

type MockDB struct {
	mock.Mock
}

func (MockDB) MustBegin() *sqlx.Tx {
	fmt.Println("In MockDB")
	db, _ := sqlx.Open("sqlite3", ":memory:")
	return db.MustBegin()
}

type MockLocationDao struct {
	mock.Mock
}

func (MockLocationDao) Insert(tx sqlx.Tx, location model.Location) int {
	return 1
}

func (MockLocationDao) GetAll(tx sqlx.Tx) []model.Location {
	return make([]model.Location, 0)
}

func (MockLocationDao) Map(tx sqlx.Tx) {
}

type MockTransitionDao struct {
	mock.Mock
}

func (MockTransitionDao) Insert(tx sqlx.Tx, locationId int, transition model.Transition) {
}

func (MockTransitionDao) GetAll(tx sqlx.Tx) []model.Transition {
	return make([]model.Transition, 0)
}

func TestService(t *testing.T) {
	binder := axon.NewBinder(
		axon.NewPackage(
			axon.Bind("LocationService").
				To().StructPtr(new(LocationService)),
		))
	injector := axon.NewInjector(binder)
	injector.Add("Db", axon.StructPtr(new(MockDB)))
	injector.Add("LocationDao", axon.StructPtr(new(MockLocationDao)))
	injector.Add("TransitionDao", axon.StructPtr(new(MockTransitionDao)))

	service := injector.GetStructPtr("LocationService").(*LocationService)
	service.GetAll()
}
