package terminal

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"golang.org/x/term"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	defaultCols  = 80
	defaultLines = 24
)

var (
	interactive bool
	nocolor     bool
	cols        int
)

func init() {
	interactive = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	// terminal columns
	cols = defaultCols
	if !IsInteractive() {
		return
	}

	c, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return
	}

	cols = c
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

func ActionNoColon(level Level, msg string) {
	caret(level)
	SetColor(color.FgWhite)
	fmt.Print(msg)
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
	if !IsInteractive() {
		if ok {
			fmt.Println(" OK")
		} else {
			fmt.Println(" FAIL")
		}

		return
	}

	if ok {
		SetColor(color.BgHiGreen, color.FgHiWhite)
		fmt.Print("  OK  ")
	} else {
		SetColor(color.BgHiRed, color.FgHiWhite)
		fmt.Print(" FAIL ")
	}

	ResetColor()
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

func DashedLine(fromCol int) {
	SetColor(color.FgHiWhite)
	fmt.Print(strings.Repeat("_", cols-fromCol-6))
	ResetColor()
}
