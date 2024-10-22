package terminal

import (
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

var (
	interactive bool
	nocolor     bool
)

func init() {
	interactive = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

func IsInteractive() bool {
	if nocolor {
		return false
	}
	return interactive
}

func SetNoColor(forceNoColor bool) {
	nocolor = forceNoColor
}

func caret(level Level) {
	switch level {
	case DebugLevel:
		SetColor(color.FgHiBlue)
	case InfoLevel:
		SetColor(color.FgHiGreen)
	case WarnLevel:
		SetColor(color.FgHiYellow)
	case ErrorLevel:
		SetColor(color.FgHiRed)
	}
	fmt.Print(">")
	SetColor(color.Reset)
	fmt.Print(" ")
}

func Line(level Level, msg string) {
	caret(level)
	fmt.Println(msg)
}

func Action(level Level, msg string) {
	caret(level)
	SetColor(color.FgWhite)
	fmt.Print(msg)
	SetColor(color.FgWhite)
	fmt.Print(": ")
	SetColor(color.Reset)
}

func Error(err error) {
	if err == nil {
		return
	}

	caret(ErrorLevel)
	SetColor(color.BgRed, color.FgWhite)
	fmt.Print("ERROR")
	SetColor(color.Reset)
	SetColor(color.FgHiWhite)
	fmt.Print(": ")
	SetColor(color.FgWhite)
	fmt.Print(err.Error())
	SetColor(color.Reset)
	fmt.Println()
}

func Result(ok bool) {
	if ok {
		SetColor(color.BgHiGreen, color.FgHiWhite)
		fmt.Print(" OK ")
	} else {
		SetColor(color.BgHiRed, color.FgHiWhite)
		fmt.Print(" FAIL ")
	}

	SetColor(color.Reset)
	fmt.Println()
}

func TableTile(title string) {
	SetColor(color.BgHiBlue, color.FgWhite)
	fmt.Print("  " + title + ":  ")
	SetColor(color.Reset)
	fmt.Println()
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
