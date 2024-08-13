package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"os/user"
	"path/filepath"
)

// fileExists returns true if the given file exists, and false if it
// doesn't. If it encounters an error, it prints the error and exits.
func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		printErrAndExit(err.Error())
		return false // NEVER REACHED
	}
}

// mustExist can be called to ensure that a file exists; it errors and exits if
// the file doesn't exist.
func mustExist(filename string) {
	if fileExists(filename) != true {
		printErrAndExit(os.ErrNotExist.Error())
	}
}

// getConfig fetches the config file name for the given file extension.
// It returns two values: the first is true if the config file exists.
// If it does, the second value is the config filename.
// If it doesn't, the second value is blank and can be ignored.
func getConfig(configPath, extension string) (bool, string) {
	if extension == "" {
		return false, ""
	}
	// Assuming the file has an extension
	fileName := filepath.Join(configPath, extension[1:]+".conf")
	if exists := fileExists(fileName); exists == false {
		return false, ""
	} else {
		return true, fileName
	}
}

// printFile is used when no config file can be found for the file extension
// It prints out the file as it reads it, with no modifications applied. Essentially
// works like 'cat'.
func printFile(fileName string) {
	mustExist(fileName)
	data, err := os.ReadFile(fileName)
	if err != nil {
		printErrAndExit(err.Error())
	}
	fmt.Print(string(data))
	return
}

// computeLineNumDigits computes the number of digits in the number of lines
// in the given byte array.
func computeLineNumDigits(data []byte) int {
	numLines := bytes.Count(data, []byte{'\n'}) + 1
	return int(math.Round(math.Log10(float64(numLines))))
}

func main() {
	disableColorFlag := flag.Bool("d", false, "Disable color")
	lineNumberFlag := flag.Bool("n", false, "Print line numbers")
	// Used only if lineNumberFlag is true
	var lineNumDigits int
	var lineNum int
	flag.Parse()

	// Check if config exists. If it doesn't, generate the config files.
	var configPath string // Location of config files, depends on OS
	if runningOnWindows() {
		configPath = "%APPDATA%\\ccat"
	} else {
		currentUser, err := user.Current()
		if err != nil {
			panic(err)
		}
		configPath = filepath.Join("/home/" + currentUser.Username + "/.config/ccat/")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		generateDefaultConfigs()
	}

	// Check if user has provided a file name
	if len(flag.Args()) < 1 {
		printErrAndExit("No File specified")
	}
	fileName := flag.Args()[0]

	// Check if file exists.
	mustExist(fileName)

	extension := filepath.Ext(fileName)
	configExists, configFilename := getConfig(configPath, extension)
	// If the given file has no corresponding config, print the file out and exit.
	if configExists == false {
		printFile(fileName)
		return
	}

	// To save computing time, determine here if the file is empty. If it is, exit
	// the program.
	finfo, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}
	if finfo.Size() == 0 {
		os.Exit(0)
	}

	// Assuming the file is not empty...
	// If the given file has a config, load the config into a stack of regColors.
	regColorStack, err := loadConfig(configFilename)
	if err != nil {
		printErrAndExit(err.Error())
	}

	// Load the input file into a colorunit slice (units) and a byte slice (data)
	units, data, err := loadInputFile(fileName)
	if err != nil {
		printErrAndExit(err.Error())
	}

	// If the '-n' flag is set, compute the number of digits in the number of lines
	// in the file, to determine the padding for the line numbers.
	if *lineNumberFlag {
		lineNumDigits = computeLineNumDigits(data)
	}
	// For each regular expression in the stack, apply it to the byte slice. Find
	// the first and last index of all matches of the regex. Then apply the corresponding color
	// to every character within these indices.
	//
	// The infinite for loop exists, because I couldn't figure out a way to pop an element from
	// the stack inside the 'for' statement. The loop exits when the 'pop' call returns 'false',
	// indicating that the stack is empty.
	//
	// The loop is also only run if the 'disable color' flag is not set.
	for *disableColorFlag == false {
		regclr, ok := regColorStack.Pop()
		// regColorStack.Pop() returns false when there are no more elements to pop
		if ok != true {
			break
		}
		re := regclr.re
		clr := regclr.clr
		// Returns an int double-slice, where each slice contains the start and end indices
		// of the match. In this case, I am finding all the matches of 're' in 'data'.
		matches := re.FindAllSubmatchIndex(data, -1)
		if matches == nil {
			continue
		}
		// For each match, apply the corresponding color to all characters in the match.
		for _, match := range matches {
			units = applyColor(units, match[0], match[1], clr)
		}
	}

	// After all possible regexes have been matched, print out the contents of 'units'.

	// If the line number flag is set, initialize the lineNum variable and print the first line number
	// with the appropriate padding.
	if *lineNumberFlag {
		lineNum = 1
		fmt.Printf("   %*d  ", lineNumDigits, lineNum)
		lineNum++
	}
	for idx, unit := range units {
		unit.print()
		// If the flag is set and we encounter a newline (and the newline isn't a trailing newline),
		// then print the next line number.
		if *lineNumberFlag && unit.ch == '\n' && idx != len(units)-1 {
			fmt.Printf("   %*d  ", lineNumDigits, lineNum)
			lineNum++
		}
	}
}
