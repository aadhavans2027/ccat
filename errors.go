package main

import (
	"fmt"
	"os"
)

func printErrAndExit(errorStr string) {
	fmt.Printf("ERROR: %s\n", errorStr)
	os.Exit(1)
}
