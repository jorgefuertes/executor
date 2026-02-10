package terminal

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jorgefuertes/executor/internal/config"
	"github.com/muesli/termenv"
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

type Term struct {
	cfg         *config.Config
	color       bool
	interactive bool
	width       int
}

func New(cfg *config.Config) *Term {
	t := &Term{
		cfg:         cfg,
		color:       lipgloss.ColorProfile() != termenv.Ascii,
		interactive: lipgloss.ColorProfile() != termenv.Ascii,
		width:       defaultCols,
	}

	if cfg.NoColor {
		t.color = false
	}

	if cfg.NoInteractive {
		t.interactive = false
	}

	t.HideCursor()
	c, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		t.interactive = false

		return t
	}

	t.width = c

	return t
}

func (t *Term) CleanUp() {
	t.ShowCursor()
}

func (t *Term) SetNoInteractive() {
	t.interactive = false
}

func (t *Term) IsInteractive() bool {
	return t.interactive
}

func (t *Term) HasColor() bool {
	return t.color
}

func (t *Term) Result(ok bool) {
	if ok {
		okText := "  OK  "
		if !t.HasColor() {
			okText = "[ OK ]"
		}
		t.Print(SuccessLabelColor, false, okText)
	} else {
		failText := " FAIL "
		if !t.HasColor() {
			failText = "[FAIL]"
		}
		t.Print(ErrorLabelColor, false, failText)
	}

	fmt.Println()
}

func (t *Term) TableTile(title string) {
	t.PrintF(TableTitleColor, false, " %s: ", title)
	fmt.Println()
}

func (t *Term) HideCursor() {
	if !t.IsInteractive() {
		return
	}

	fmt.Print("\033[?25l")
}

func (t *Term) ShowCursor() {
	if !t.IsInteractive() {
		return
	}

	fmt.Print("\033[?25h")
}

func (t *Term) DashedLine() {
	if !t.IsInteractive() {
		t.Print(SecondaryColor, false, ellipsis)

		return
	}

	_, col, err := t.GetCursorPosition()
	if err != nil {
		t.Print(SecondaryColor, false, ellipsis)

		return
	}

	t.Print(SecondaryColor, false, strings.Repeat(ellipsis, t.width-col))
}

func (t *Term) GetCursorPosition() (int, int, error) {
	if !t.IsInteractive() {
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
