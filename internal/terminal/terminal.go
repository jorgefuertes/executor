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
	defaultCols                 = 80
	defaultLines                = 24
	minInteractiveRemainingCols = 15
)

var (
	interactive bool
	nocolor     bool
	width       int
)

func init() {
	interactive = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	// terminal columns
	width = defaultCols
	if !IsInteractive() {
		SetNoColor(true)

		return
	}

	HideCursor()
	c, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		interactive = false

		return
	}

	width = c
}

func CleanUp() {
	ShowCursor()
}

func IsInteractive() bool {
	return interactive
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

func DashedLine() {
	if !IsInteractive() {
		Print(SecondaryColor, false, ellipsis)

		return
	}

	_, col, err := GetCursorPosition()
	if err != nil {
		Print(SecondaryColor, false, ellipsis)

		return
	}

	Print(SecondaryColor, false, strings.Repeat(ellipsis, width-col))
}

func GetCursorPosition() (int, int, error) {
	if !IsInteractive() {
		return 0, 0, fmt.Errorf("not an interactive terminal")
	}

	// save current terminal state
	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get terminal state: %w", err)
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
	}()

	// set terminal to raw mode
	_, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to set raw mode: %w", err)
	}

	// request cursor position
	fmt.Print("\033[6n")

	// read response with timeout
	buf := make([]byte, 32)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read cursor position: %w", err)
	}

	// Parse response: \033[{row};{col}R
	var row, col int
	_, err = fmt.Sscanf(string(buf[:n]), "\033[%d;%dR", &row, &col)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse cursor position: %w", err)
	}

	return row, col, nil
}
