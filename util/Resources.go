package util

import (
	"fmt"
	"os"
)

func Resource(path string) ([]byte, error) {
	pwd, _ := os.Getwd()
	dat, err := os.ReadFile(fmt.Sprintf("%s/%s", pwd, path))
	if (err != nil) {
		fmt.Printf("Crash: %s", err)
		os.Exit(1)
	}
	return dat, err
}

