package main

import (
	"takeoff.com/matilda/applications"
)

func main() {
	var app = applications.InitApplication()
	app.LocationService.Display()
}
