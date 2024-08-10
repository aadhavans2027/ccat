package main

import (
	"fmt"

	colorData "github.com/fatih/color"
)

// A color represents a possible color, which text can be printed out in.
// Each color has a name and an object (from fatih/color). This object is used
// to print text in that color.
type color struct {
	name     string
	colorObj *colorData.Color
}

// The following is a list of all possible colors, stored in a map.
var possibleColors map[string]color = map[string]color{
	"BLACK":   {"BLACK", colorData.New(colorData.FgBlack)},
	"RED":     {"RED", colorData.New(colorData.FgRed)},
	"GREEN":   {"GREEN", colorData.New(colorData.FgGreen)},
	"YELLOW":  {"YELLOW", colorData.New(colorData.FgYellow)},
	"BLUE":    {"BLUE", colorData.New(colorData.FgBlue)},
	"MAGENTA": {"MAGENTA", colorData.New(38, 2, 254, 141, 255)},
	"CYAN":    {"CYAN", colorData.New(colorData.FgCyan)},
	"WHITE":   {"WHITE", colorData.New(colorData.FgWhite)},
	"GRAY":    {"GRAY", colorData.New(colorData.FgWhite, colorData.Faint)},
	// Last three numbers are RGB. Reference https://en.wikipedia.org/wiki/ANSI_escape_code for what the first two numbers mean.
	"ORANGE": {"ORANGE", colorData.New(38, 2, 255, 153, 28)},
	"NONE":   {"NONE", colorData.New()},
}

// Apply the given color 'clr' to all units in 'units', within the indices
// marked by 'start' and 'end'
func applyColor(units []colorunit, start int, end int, clr color) []colorunit {
	for i := start; i < end; i++ {
		units[i].clr = clr
	}
	return units
}

// newColor takes a string, and if it represents one of the colors in the dictionary,
// it returns the appropriate color. If it doesn't, the function returns an error.
func newColor(colorString string) (color, error) {
	clr, ok := possibleColors[colorString]
	if ok != true {
		return color{}, fmt.Errorf("Invalid color: %s", colorString)
	}
	return clr, nil
}

// newColorMust is similar to newColor, but panics if the given color isn't valid.
func newColorMust(colorString string) color {
	if clr, err := newColor(colorString); err != nil {
		panic(err)
	} else {
		return clr
	}
}
