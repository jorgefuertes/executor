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
		SetNoColor(true)
		return
	}

	HideCursor()
	c, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return
	}

	cols = c
}

func CleanUp() {
	color.Unset()
	ShowCursor()
}

func IsInteractive() bool {
	if nocolor {
		return false
	}
	return interactive
}

func SetNoColor(forceNoColor bool) {
	nocolor = forceNoColor
	color.NoColor = nocolor
}

func Line(level Level, msg string, slow bool) {
	caret(level)
	Print(White, slow, msg+"\n")
}

func Action(level Level, msg string, slow bool) {
	caret(level)
	Print(White, slow, msg+": ")
}

func Error(err error) {
	if err == nil {
		return
	}

	caret(ErrorLevel)
	Print(RedLabel, false, "ERROR")
	PrintF(White, false, ": %s", err.Error())
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
		Print(GreenLabel, false, "  OK  ")
	} else {
		Print(RedLabel, false, " FAIL ")
	}

	fmt.Println()
}

func TableTile(title string) {
	PrintF(CyanLabel, false, " %s: ", title)
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
	Print(Gray, false, strings.Repeat("_", cols-fromCol-6))
}
