/*
TODO:
    Install sqllite: https://flaviocopes.com/sqlite-how-to-install
	Get GORM working with sqllite for now: https://gorm.io/docs/index.html
	Clear and populate our database
	Load a board and create a graph of adjacent nodes
		Static items like posts are unmovable
    	Shelves are a static item that can contain slots
		Within slots we can store bins or items
	We will be tracking the pickers, pallets and the items
		For now let's assume pickers all walk at the same pass
		Using the graph above we should be able to calculate a path from the picker to target
		Record the location of our pickers or items locations when a job is started or finished
		If a picker finishes a job they should request a new job and have one presented based on location.
		visit our locations like a graph, avoiding cycles record all paths between points upfront
		in comparing paths we can anticipate collisions but don't get hung up on that because they are humans
		create a random field which changes the offsets of out pickers
	A list of products and there location: on a pallet, on a shelf, or in a bin.
	Now create a list of things to pick and see if we can arrange there jobs efficiently.
	Try to build a display using:
		https://github.com/jkomoros/boardgame
		https://github.com/jkomoros/boardgame/blob/master/TUTORIAL.md
	https://github.com/eddieowens/axon
	Get this all to work on the cloud.
	Good article on tags: https://medium.com/golangspec/tags-in-golang-3e5db0b8ef3e
	Magical comments: https://blog.jbowen.dev/2019/09/the-magic-of-go-comments
	Simpler JSON? https://pkg.go.dev/github.com/bitly/go-simplejson
	Break project into: src, data/resources, ... it currently looks weird to have our json in the middle of go files
	interfaces and delegation: https://medium.com/code-zen/go-interfaces-and-delegation-pattern-f962c138dc1e
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
	"matilda/floor"
	"os"
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
