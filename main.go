package main

import (
	"takeoff.com/matilda/applications"
)

func main() {
	app := applications.Get()
	app.LocationService.GetAll()
}
