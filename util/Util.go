package util

import (
	"fmt"
	"os"
	"takeoff.com/matilda/data"
)

func LoadPlan(path string) ([]byte, error) {
	dat, err := os.ReadFile(data.Path(path))
	if err != nil {
		fmt.Printf("Crash: %s", err)
		os.Exit(1)
	}
	return dat, err
}
