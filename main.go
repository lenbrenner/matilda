package main

import (
	"takeoff.com/matilda/applications"
)

func main() {
	applications.Get().LocationService.Display()
}
