package util

import (
	"fmt"
	"os"
)

func ResourcePath(path string) string {
	pwd, _ := os.Getwd()
	return fmt.Sprintf("%s/%s", pwd, path)
}

func Resource(path string) ([]byte, error) {
	dat, err := os.ReadFile(ResourcePath(path))
	if err != nil {
		fmt.Printf("Crash: %s", err)
		os.Exit(1)
	}
	return dat, err
}
