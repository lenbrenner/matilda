/*
*/

//Create some locations
//Create some products
//Create some people
//Using locations build paths
//Using paths find intersections.
//This is how slotting is optimized: https://github.com/TakeoffTech/SDIA#sdia-solver-for-dynamic-inventory-allocation-former-continuous-slotting
package main

import (
	"fmt"
	"os"
	"takeoff.com/matilda/floor"
)

func pathsTo(location floor.Location, destination floor.Location) {
	//fmt.Printf("Searching for paths from %v to %v")
	//seen := map[LocationLabel]bool{}
	for direction, label := range location.Transitions {
		//node := locations[label]
		fmt.Println("%v, %v", direction, label)

	}
}

type locationArray []floor.Location

func main() {
	//gormExample()
	//barr, _ := json.MarshalIndent(locations, "", "    ")
	//fmt.Println(string(barr))
	var plan floor.Plan
	err := plan.LoadJson("3x3_floor")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, location := range plan.Locations {
		fmt.Printf("%v %v\n", i, location)
		//for j, destination := range locations {
		//	if i!=j {
		//		location.pathsTo(destination)
		//	}
		//}
	}
}

/*
package main

func main() {
g := &Game{Health: 100, Welcome: "Welcome to the Starship Enterprise\n\n", CurrentLocation: "Bridge"}
g.Play()
}
*/
