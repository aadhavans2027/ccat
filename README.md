## ccat

ccat is a file printing tool (like 'cat') which uses Regular Expressions to enable syntax highlighting.

---

### Features
- Support for 11 colors: Red, Blue, Green, Magenta, Cyan, Black, White, Yellow, Gray, Orange and Dark Blue.
- Adding more colors involves adding a line of code, then recompiling.
- Regex-color mappings are stored in configuration files.
- Uses the file extension to determine which configuration file to use.
- Highly extensible - to add a config file for an specific file type, name the file `<extension>.conf`.
- Support for printing line numbers with the `-n` flag.
- Statically linked Go binary - no runtime dependencies, config files are distributed along with the binary.
- Cross-platform

---

### Installing
If you have the `go` command installed, run `make` after cloning the repository.

---

### Getting Started
The config files are embedded within the binary. They will automatically be installed to the correct location (`%APPDATA/ccat` on Windows, `~/.config/ccat` on UNIX) when the program is first run.

As written above, if provided a file with extension `.example`, the program will look for the config file named `example.conf`. If such a file doesn't exist, the file is printed out without any highlighting.

---

### Config Files

The config files are written in YAML, and have the following syntax:

`"<regex>": COLOR`

Note that the regex must be enclosed in double quotes, and the color must be capitalized.

---

### TODO:
- Allow user to define colors at runtime by reading RGB values from a config file.
- Allow users to provide a config file in the command-line, overriding the extension-based config file.
- Provide releases.
