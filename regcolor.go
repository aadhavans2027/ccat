package main

import "gitea.twomorecents.org/Rockingcool/kleingrep/regex"

// A regColor is a regex-color pair. The config file is read
// into a stack of this data type.
type regColor struct {
	re  *regex.Reg
	clr color
}
