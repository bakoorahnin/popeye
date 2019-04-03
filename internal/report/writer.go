package report

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/derailed/popeye/internal/linter"
)

// Color ANSI palette (256!)
const (
	ColorOrangish  = 220
	ColorOrange    = 208
	ColorGray      = 250
	ColorWhite     = 15
	ColorBlue      = 105
	ColorRed       = 202
	ColorCoolBlue  = 99
	ColorAqua      = 122
	ColorDarkOlive = 155
	ColorLighSlate = 75 // 105
)

// FontBold style
const (
	FontBold = 1
)

// Color tracks the output color.
type Color int

const (
	reportWidth = 80
	tabSize     = 2
)

// Open begins a new report section.
func Open(w io.Writer, s string) {
	fmt.Fprintf(w, "\n%s\n", Colorize(s, ColorLighSlate))
	fmt.Fprintf(w, "%s\n", Colorize(strings.Repeat("┅", 80), ColorLighSlate))
	fmt.Fprintln(w)
}

// Close a report section.
func Close(w io.Writer) {
	fmt.Fprintln(w)
}

// Error prints out error out.
func Error(w io.Writer, fmat string, args ...interface{}) {
	fmt.Fprintln(w)
	msg := fmt.Sprintf(fmat, args...)
	buff := make([]string, 0, len(msg)%reportWidth)
	width := reportWidth - 3
	for i := 0; len(msg) > width; i += width {
		buff = append(buff, msg[i:i+width])
		msg = msg[i+width:]
	}
	buff = append(buff, msg)
	fmt.Fprintf(w, "💥 "+Colorize(strings.Join(buff, "\n"), ColorRed))
	fmt.Fprintln(w)
}

// Comment writes a comment line.
func Comment(w io.Writer, msg string) {
	fmt.Fprintf(w, "  · "+msg+"\n")
}

// Dump all errors to output.
func Dump(w io.Writer, l linter.Level, issues ...linter.Issue) {
	var current string
	for _, i := range issues {
		if i.Severity() >= l {
			tokens := strings.Split(i.Description(), linter.Delimiter)
			if len(tokens) == 1 {
				Write(w, i.Severity(), 2, i.Description()+".")
			} else {
				if current != tokens[0] {
					Write(w, containerLevel, 2, tokens[0])
					current = tokens[0]
				}
				Write(w, i.Severity(), 3, tokens[1]+".")
			}
		}
	}
}

// Write a colorized message to stdout.
func Write(w io.Writer, l linter.Level, indent int, msg string) {
	spacer := strings.Repeat(" ", tabSize*indent)

	if indent == 1 {
		dots := reportWidth - len(msg) - tabSize*indent - 3
		msg = Colorize(msg, colorForLevel(l)) + Colorize(strings.Repeat(".", dots), ColorGray)
		fmt.Fprintf(w, "%s· %s%s\n", spacer, msg, emojiForLevel(l))
		return
	}

	msg = Colorize(msg, ColorWhite)
	fmt.Fprintf(w, "%s%s %s\n", spacer, emojiForLevel(l), msg)
}

// Colorize a string based on given color.
func Colorize(s string, c Color) string {
	return "\033[38;5;" + strconv.Itoa(int(c)) + ";m" + s + "\033[0m"
}

func colorForLevel(l linter.Level) Color {
	switch l {
	case linter.ErrorLevel:
		return ColorRed
	case linter.WarnLevel:
		return ColorOrangish
	case linter.InfoLevel:
		return ColorAqua
	default:
		return ColorAqua
	}
}

const containerLevel linter.Level = 100

func emojiForLevel(l linter.Level) string {
	switch l {
	case containerLevel:
		return emojis["container"]
	case linter.ErrorLevel:
		return emojis["farfromfok"]
	case linter.WarnLevel:
		return emojis["warn"]
	case linter.InfoLevel:
		return emojis["fyi"]
	default:
		return emojis["peachy"]
	}
}

// Logo popeye
var Logo = []string{
	"K          .-'-.     ",
	" 8     __|      `\\  ",
	"  s   `-,-`--._   `\\",
	" []  .->'  a     `|-'",
	"  `=/ (__/_       /  ",
	"    \\_,    `    _)  ",
	"       `----;  |     ",
}

// Popeye title
var Popeye = []string{
	` ___     ___ _____   _____ `,
	`| _ \___| _ \ __\ \ / / __|`,
	`|  _/ _ \  _/ _| \ V /| _| `,
	`|_| \___/_| |___| |_| |___|`,
}