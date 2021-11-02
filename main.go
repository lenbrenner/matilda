package main

import (
	"takeoff.com/matilda/applications"
)

//Todo this needs implementing.
//func pathsTo(location model.Location, destination model.Location) {
//	//fmt.Printf("Searching for paths from %v to %v")
//	//seen := map[LocationLabel]bool{}
//	for direction, label := range location.Transitions {
//		//node := locations[label]
//		fmt.Println("%v, %v", direction, label)
//
//	}
//}

func main() {
	var app = applications.InitApplication()
	app.Service.DoSomething()
}
