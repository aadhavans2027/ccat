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

### Getting Started
The config files are embedded within the binary. They will automatically be installed to the correct location (`%APPDATA/ccat` on Windows, `~/.config/ccat` on UNIX) when the program is first run.
TODO:
- Allow user to define colors at runtime by reading RGB values from a config file.
- Provide releases.
