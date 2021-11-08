package main

import (
	"fmt"
	"takeoff.com/matilda/applications"
)

func main() {
	app := applications.Get()
	for _, x := range app.LocationService.GetAll() {
		 fmt.Println(x)
	}
}
