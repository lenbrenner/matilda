package applications

import (
	"fmt"
	"github.com/eddieowens/axon"
	"github.com/stretchr/testify/mock"
	"takeoff.com/matilda/model"
	"testing"
)

type MockLocationService struct {
	mock.Mock
}

func (MockLocationService) LoadAll(locations []model.Location) {
	fmt.Println("In LocationService")
}

func (MockLocationService) GetAll() []model.Location {
	fmt.Println("In MockLocationService.LoadAll")
	return make([]model.Location, 0)
}

func TestApplication(t *testing.T) {
	binder := axon.NewBinder(
		axon.NewPackage(
			axon.Bind("Application").
				To().StructPtr(new(Application)),
		))
	injector := axon.NewInjector(binder)
	injector.Add("LocationService", axon.StructPtr(new(MockLocationService)))

	application := injector.GetStructPtr("Application").(*Application)
	application.LoadPlan("3x3_floor")
}