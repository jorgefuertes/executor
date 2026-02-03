package terminal

import (
	"fmt"
	"os"
	"strings"

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
	ShowCursor()
}

func IsInteractive() bool {
	return interactive
}

func Line(level Level, msg string, slow bool) {
	caret(level)
	Print(PrimaryColor, slow, msg+"\n")
}

func Action(level Level, msg string, slow bool) {
	caret(level)
	Print(PrimaryColor, slow, msg+": ")
}

func Error(err error) {
	if err == nil {
		return
	}

	caret(ErrorLevel)
	Print(ErrorColor, false, "ERROR")
	PrintF(PrimaryColor, false, ": %s", err.Error())
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
		okText := "  OK  "
		if !HasColor() {
			okText = "[ OK ]"
		}
		Print(SuccessLabelColor, false, okText)
	} else {
		failText := " FAIL "
		if !HasColor() {
			failText = "[FAIL]"
		}
		Print(ErrorLabelColor, false, failText)
	}

	fmt.Println()
}

func TableTile(title string) {
	PrintF(TableTitleColor, false, " %s: ", title)
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

func DashedLine(fromCol int, rightMargin int) {
	repeat := cols - fromCol - rightMargin
	if repeat < 0 {
		repeat = 0
	}

	Print(SecondaryColor, false, strings.Repeat("…", repeat))
}
