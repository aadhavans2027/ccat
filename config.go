package main

import (
	"ccat/stack"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// loadConfig takes in the filename of a config file. It reads the file,
// and returns a stack of RegColors, with the item at the bottom being the one that
// was read first. This ensures that, _when accessing the RegColors in the stack, the last
// one (ie. the one that was read first) has highest precedence_.
// If there is an error compiling the regular expressions, the error is returned.
func loadConfig(configFilename string) (stack.Stack[regColor], error) {
	configFile, err := os.ReadFile(configFilename)
	if err != nil {
		return *stack.NewStack[regColor](0), err
	}

	// Here, I create a MapSlice. This is a slice of key-value pairs, and will
	// store the results of unmarshalling the YAML file.
	tempMapSlice := yaml.MapSlice{}
	if err := yaml.Unmarshal(configFile, &tempMapSlice); err != nil {
		return *stack.NewStack[regColor](0), err
	}

	// Here, I create the stack which will eventually be returned.
	// Each element of the MapSlice (created above) stores the key and value of a line
	// in the file.
	// Each regex string is compiled, and if there is an error, that error is
	//  returned.
	regColorStack := stack.NewStack[regColor](len(strings.Split(string(configFile), "\n"))) // The stack will have the same size as the number of lines in the file
	for _, item := range tempMapSlice {
		re := regexp.MustCompile(item.Key.(string))
		clr, err := newColor(item.Value.(string))
		if err != nil {
			return *stack.NewStack[regColor](0), err
		}
		// If we got past the errors, then the color _must_ be valid.
		regColorStack.Push(regColor{re, clr})
	}

	return *regColorStack, nil
}
