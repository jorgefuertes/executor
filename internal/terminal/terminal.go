package terminal

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

func caret(level Level) {
	switch level {
	case DebugLevel:
		color.Set(color.FgHiBlue)
	case InfoLevel:
		color.Set(color.FgHiGreen)
	case WarnLevel:
		color.Set(color.FgHiYellow)
	case ErrorLevel:
		color.Set(color.FgHiRed)
	}
	fmt.Print(">")
	color.Set(color.Reset)
	fmt.Print(" ")
}

func Line(level Level, msg string) {
	caret(level)
	fmt.Println(msg)
}

func Action(level Level, msg string) {
	caret(level)
	color.Set(color.FgWhite)
	fmt.Print(msg)
	color.Set(color.FgWhite)
	fmt.Print(": ")
	color.Set(color.Reset)
}

func Error(err error) {
	if err == nil {
		return
	}

	caret(ErrorLevel)
	color.Set(color.BgRed, color.FgWhite)
	fmt.Print("ERROR")
	color.Set(color.Reset)
	color.Set(color.FgHiWhite)
	fmt.Print(": ")
	color.Set(color.FgWhite)
	fmt.Print(err.Error())
	color.Set(color.Reset)
	fmt.Println()
}

func ActionError(err error) {
	Result(err == nil)
	Error(err)
}

func Result(ok bool) {
	if ok {
		color.Set(color.BgHiGreen, color.FgHiWhite)
		fmt.Print(" OK ")
	} else {
		color.Set(color.BgHiRed, color.FgHiWhite)
		fmt.Print(" FAIL ")
	}

	color.Set(color.Reset)
	fmt.Println()
}

func TableTile(title string) {
	color.Set(color.BgHiBlue, color.FgWhite)
	fmt.Print("  " + title + ":  ")
	color.Set(color.Reset)
	fmt.Println()
}

func ShowOutput(level Level, title string, out bytes.Buffer) {
	if out.Len() == 0 {
		fmt.Println()
		Line(level, title+": Empty")
		fmt.Println()

		return
	}

	fmt.Println()
	Line(level, "---[BEGIN "+title+"]---")
	fmt.Println()
	fmt.Println(out.String())
	Line(level, "---[END "+title+"]----")
}

func IsInteractive() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

func SavePos() {
	if !IsInteractive() {
		return
	}

	fmt.Print("\033[s")
}

func RestorePos() {
	if !IsInteractive() {
		return
	}

	fmt.Print("\033[u")
}

func HideCursor() {
	if !IsInteractive() {
		return
	}

	fmt.Print("\033[?25l")
}

func ShowCursor() {
	if !IsInteractive() {
		return
	}

	fmt.Print("\033[?25h")
}
