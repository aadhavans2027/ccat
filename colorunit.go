package main

import (
	"os"
)

// A colorunit represents a unit in a file. It consists of the character,
// and the color that the character should be printed out in.
type colorunit struct {
	ch  byte
	clr color
}

// loadInputFile loads the given file and returns a slice of colorunits,
// and a slice of bytes (which just contains all the text in the file).
// The slice of colorunits is used to fill in the color for each character.
// The slice of bytes is used to perform the regex matching.
// The color will be set to the current terminal foreground color.
func loadInputFile(fileName string) ([]colorunit, []byte) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	units := make([]colorunit, len(data))
	for idx, c := range data {
		units[idx] = colorunit{byte(c), newColorMust("NONE")}
	}
	return units, data
}

// print is used to print out the character in the given colorunit, according to
// its color.
func (unit colorunit) print() {
	unit.clr.colorObj.Printf("%c", unit.ch)
	return
}
