package main

import (
	"embed"
	"errors"
	"gitea.twomorecents.org/Rockingcool/ccat/stack"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

//go:embed config
var storedConfigs embed.FS // Embed the folder containing config files

// runningOnWindows: At the moment this function isn't used. When Window support is added,
// it will be used to determine if the program is being run on Windows.
func runningOnWindows() bool {
	return runtime.GOOS == "windows"
}

// generateDefaultConfigs is used to generate a folder of default config files
// for common languages. These default config files are embedded into the program, and will
// be outputted into the given directory.
//
// If there is an error encountered, the error is returned.
func generateDefaultConfigs(configOutputPath string) error {
	err := os.MkdirAll(configOutputPath, 0755)
	if err != nil {
		if os.IsExist(err) {
			return errors.New("Directory already exists.")
		} else {
			return errors.New("Unable to create directory.")
		}
	}

	// Copy each folder from the embedded filesystem, into the destination path
	err = fs.WalkDir(storedConfigs, "config", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() { // Skip directories
			return nil
		}
		relPath, _ := filepath.Rel("config", path)
		dstPath := filepath.Join(configOutputPath, relPath) // Destination path

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err := os.WriteFile(dstPath, data, 0644); err != nil {
			return err
		}
		return nil
	})
	return nil
}

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
