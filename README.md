## ccat

ccat is a file printing tool (like 'cat') which uses Regular Expressions to enable syntax highlighting.

---

### Features
- 11 colors are defined out-of-the-box: RED, BLUE, GREEN, MAGENTA, CYAN, BLACK, WHITE, YELLOW, GRAY, ORANGE and DARKBLUE.
- Support for defining custom colors via the `ccat.colors` file.
- Regex-color mappings are stored in configuration files.
- Uses the file extension to determine which configuration file to use.
- Highly extensible - to add a config file for an specific file type, name the file `<extension>.conf`.
- Support for printing line numbers with the `-n` flag.
- Statically linked Go binary - no runtime dependencies, config files are distributed along with the binary.
- Linux and MacOS supported.

---

### Installing

Download the appropriate zip-file from the 'Releases' section. Place the executable in your PATH.

NOTE: The releases are not available on the GitHub repo (which is a mirror of https://gitea.twomorecents.org/Rockingcool/ccat). Obtain the [releases](https://gitea.twomorecents.org/Rockingcool/ccat/releases) from there instead.

---

### Building from source

If you have the `go` command installed, run `make` after cloning the repository.

---

### Supported Languages

The following languages have config files included by default:

- C
- Go

---

### Getting Started
The config files are embedded within the binary. They will automatically be installed to the correct location (`~/.config/ccat` on UNIX) when the program is first run.

As written above, if provided a file with extension `.example`, the program will look for the config file named `example.conf`. If such a file doesn't exist, the file is printed out without any highlighting.

For example, if you want to create syntax highlighting for Java, create a file named `java.conf` in your config directory. In this file, include regular-expressions for each of the langauges's keywords, and provide a corresponding color. Use the provided `c.conf` and `go.conf` files as a starting point.

---

### Config Files

The config files are written in YAML. Each line has the following syntax:

`"<regex>": COLOR`

Note that the regex must be enclosed in double quotes, and the color must be capitalized.

---

### Custom Colors

To define a color of your own, create a file named `ccat.colors` in the config directory (mentioned above). The syntax of this file is the following:

`COLOR: <red> <green> <blue>`

Note that the color name must be capitalized (and shouldn't contain spaces). The RGB values must each be from 0 to 255.

---

### TODO:
- Windows support.
- Allow users to provide a config file in the command-line, overriding the extension-based config file.
