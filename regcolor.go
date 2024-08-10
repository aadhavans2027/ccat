package main

import "regexp"

// A regColor is a regex-color pair. The config file is read
// into a stack of this data type.
type regColor struct {
	re  *regexp.Regexp
	clr color
}
