package logging

import (
	"fmt"
	"os"

	"github.com/jwalton/gchalk"
)

var hint = gchalk.WithBrightBlack().WithItalic()

var info = gchalk.WithBrightGreen().WithBold()

var warning = gchalk.WithYellow().WithBold()

var err = gchalk.WithRed().WithBold()

var bold = gchalk.WithBold()

var italic = gchalk.WithItalic()

func PaintHint(message string) string {
	return hint.Paint(message)
}

func PaintInfo(message string) string {
	return info.Paint(message)
}
func PaintWarning(message string) string {
	return warning.Paint(message)
}

func PaintError(message string) string {
	return err.Paint(message)
}

func PaintBold(message string) string {
	return bold.Paint(message)
}

func PaintItalic(message string) string {
	return italic.Paint(message)
}

// Print an info message.
func PrintInfo(message string) {
	fmt.Println(info.Paint("INFO:", message))
}

// Print a warning message.
func PrintWarn(message string) {
	fmt.Println(warning.Paint("WARNING:", message))
}

// Print a warning message to stderr, so it does not mix with output meant for other programs.
func PrintWarnStderr(message string) {
	fmt.Fprintln(os.Stderr, warning.Paint("WARNING:", message))
}

// Print an error message.
func PrintError(message string) {
	fmt.Println(err.Paint("ERROR:", message))
}
