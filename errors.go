package main

import (
	"fmt"
	"os"
)

func printAndExit(errorStr string) {
	fmt.Printf("ERROR: %s\n", errorStr)
	os.Exit(1)
}
