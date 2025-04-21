package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	colorData "github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

// A color represents a possible color, which text can be printed out in.
// Each color has a name and an object (from fatih/color). This object is used
// to print text in that color.
type color struct {
	name     string
	colorObj *colorData.Color
}

// A RGB represents a Red, Blue, Green trio of values, along with SGR parameters.
// Each value is represented as an int. For info on SGR parameters, see:
// https://en.wikipedia.org/wiki/ANSI_escape_code#Select_Graphic_Rendition_parameters
// If 'red', 'green' and 'blue' are all -1, then the default terminal color is used.
// If some (but not all) of them are -1, an error is thrown.
type RGB struct {
	sgr1  int
	red   int
	blue  int
	green int
	sgr2  int
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
	"ORANGE":   {"ORANGE", colorData.New(38, 2, 255, 153, 28)},
	"DARKBLUE": {"DARKBLUE", colorData.New(38, 2, 0, 112, 255)},
	"NONE":     {"NONE", colorData.New()},
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

// newColorMust is similar to newColor, but prints an error and exits if the given color isn't valid.
func newColorMust(colorString string) color {
	if clr, err := newColor(colorString); err != nil {
		printErrAndExit(err.Error())
		panic(err) // NEVER REACHED
	} else {
		return clr
	}
}

// isValidColorName returns true if the given string only contains uppercase alphabetic
// characters.
func isValidColorName(colorName string) bool {
	for _, ch := range colorName {
		if (ch > 'Z' || ch < 'A') && (ch != '_') {
			return false
		}
	}
	return true
}

// stringToRGB takes a string representing an RGB five-tuple. It constructs and RGB type and
// returns it. Any errors encountered are returned. If an error is returned, it is safe to
// assume that the string doesn't represent an RGB five-tuple.
func stringToRGB(rgbString string) (*RGB, error) {
	values := strings.Split(rgbString, " ")
	// There must be three space-separated strings.
	if len(values) != 5 {
		// TODO: Instead of ignoring these errors and returning a generic error (as I do in the
		// callee), wrap the error returned from this function, inside the error returned by the callee.
		return nil, fmt.Errorf("Error parsing RGB five-tuple.")
	}
	// If any of the strings doesn't represent an integer (or is out of bounds), return an error.
	// WARNING: LAZY CODE INCOMING
	var toReturn RGB
	var err error
	toReturn.sgr1, err = strconv.Atoi(values[0])
	if err != nil {
		return nil, fmt.Errorf("Error parsing SGR1 integer: Invalid value.")
	}
	if toReturn.sgr1 < 0 || toReturn.sgr1 > 107 { // Maximum value for SGR values
		return nil, fmt.Errorf("Error parsing SGR1 integer: Out-of-bounds.")
	}
	toReturn.red, err = strconv.Atoi(values[1])
	if err != nil {
		return nil, fmt.Errorf("Error parsing RED integer: Invalid value.")
	}
	if toReturn.red < -1 || toReturn.red > 255 {
		return nil, fmt.Errorf("Error parsing RED integer: Out-of-bounds.")
	}
	toReturn.blue, err = strconv.Atoi(values[2])
	if err != nil {
		return nil, fmt.Errorf("Error parsing BLUE integer: Invalid value.")
	}
	if toReturn.blue < -1 || toReturn.blue > 255 {
		return nil, fmt.Errorf("Error parsing BLUE integer: Out-of-bounds.")
	}
	toReturn.green, err = strconv.Atoi(values[3])
	if err != nil {
		return nil, fmt.Errorf("Error parsing GREEN integer: Invalid value.")
	}
	if toReturn.green < -1 || toReturn.green > 255 {
		return nil, fmt.Errorf("Error parsing GREEN integer: Out-of-bounds.")
	}
	toReturn.sgr2, err = strconv.Atoi(values[4])
	if err != nil {
		return nil, fmt.Errorf("Error parsing SGR2 integer: Invalid value.")
	}
	if toReturn.sgr2 < 0 || toReturn.sgr2 > 107 {
		return nil, fmt.Errorf("Error parsing SGR2 integer: Out-of-bounds.")
	}

	if !(toReturn.red > 0 && toReturn.blue > 0 && toReturn.green > 0) &&
		!(toReturn.red == -1 && toReturn.green == -1 && toReturn.blue == -1) {
		return nil, fmt.Errorf("Error parsing color: All values must be positive or -1 for default terminal color.")
	}
	return &toReturn, nil
}

// loadColorsFromFile loads the colors defined in the given config file, and adds them to
// the possibleColors map. This allows the user to define custom colors at run-time.
// The colors config file has the following syntax:
// COLOR: <SGR1> <RED> <GREEN> <BLUE> <SGR2>
//
// Note that the color must be capitalized (and not contain spaces), and the R, G and B
// values must be from -1 to 255 (-1 refers to the default terminal color, and all three values
// must be -1 for this to work).
func loadColorsFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	// Read color config file into a MapSlice
	tempMapSlice := yaml.MapSlice{}
	if err := yaml.Unmarshal(data, &tempMapSlice); err != nil {
		return fmt.Errorf("Unable to read color config file: %s", filepath)
	}

	for _, item := range tempMapSlice {
		if !(isValidColorName(item.Key.(string))) {
			return fmt.Errorf("Invalid color name: %s", item.Key.(string))
		}
		var rgb *RGB
		if rgb, err = stringToRGB(item.Value.(string)); err != nil {
			return fmt.Errorf("Invalid RGB trio: %s", item.Value.(string))
		}
		// If we haven't returned an error yet, the color must be valid.
		// Add it to the map. colorData.New() expects values of type colorData.Attribute,
		// so we must cast our RGB values accordingly.
		// First, check if one of the color values is -1. If it is, they must all be negative (based
		// on the check in 'stringToRGB()'). If this is the case, don't put the color values.
		if rgb.red == -1 {
			possibleColors[item.Key.(string)] = color{
				item.Key.(string),
				colorData.New(
					colorData.Attribute(rgb.sgr2),
				),
			}
		} else {
			possibleColors[item.Key.(string)] = color{
				item.Key.(string),
				colorData.New(
					colorData.Attribute(rgb.sgr1),
					2,
					colorData.Attribute(rgb.red),
					colorData.Attribute(rgb.blue),
					colorData.Attribute(rgb.green),
					colorData.Attribute(rgb.sgr2),
				),
			}
		}
	}

	return nil
}
