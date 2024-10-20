package terminal

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/color"
)

const spinnerDelayMilliseconds = 100

type Progress struct {
	spin       int
	start      time.Time
	ctx        context.Context
	cancel     context.CancelFunc
	chars      []string
	printedLen int
}

var spinners = map[string][]string{
	"dots":        {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧"},
	"arrow":       {"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
	"star":        {"★", "☆", "★", "☆", "★", "☆", "★", "☆"},
	"circle":      {"◐", "◓", "◑", "◒", "◐", "◓", "◑", "◒"},
	"square":      {"▖", "▘", "▝", "▗", "▖", "▘", "▝", "▗"},
	"square-star": {"▌", "▞", "▛", "▙", "▟", "█", "▐", "█"},
	"line":        {"⎺", "⎻", "⎼", "⎽", "⎼", "⎻", "⎺", "⎼"},
	"line-star":   {"⎸", "⎹", "⎺", "⎻", "⎼", "⎽", "⎾", "⎿"},
	"bar":         {`|`, `/`, `-`, `\`, `|`, `/`, `-`, `\`},
	"o":           {".", "o", "O", "0", "O", "o", ".", " "},
}

func NewProgress(style string) *Progress {
	ctx, cancel := context.WithCancel(context.Background())

	chars, ok := spinners[style]
	if !ok {
		chars = spinners["dots"]
	}

	return &Progress{
		spin:       0,
		start:      time.Now(),
		ctx:        ctx,
		cancel:     cancel,
		chars:      chars,
		printedLen: 0,
	}
}

func (p *Progress) elapsed() string {
	t := time.Since(p.start)
	m := int(math.Floor(t.Minutes()))
	s := int(math.Floor(t.Seconds())) % 60
	out := "["
	if m > 0 {
		out += fmt.Sprintf("%d min", m)
	}
	if s > 0 || m == 0 {
		if m > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%d sec", s)
	}
	out += "]"

	return out
}

func (p *Progress) print() {

	if !IsInteractive() {
		return
	}

	SavePos()
	et := p.elapsed()
	spin := p.chars[p.spin]
	p.printedLen = len(et + " " + spin)
	SetColor(color.FgCyan)
	fmt.Print(et + " ")
	SetColor(color.FgYellow)
	fmt.Print(spin)
	ResetColor()
	RestorePos()
}

func (p *Progress) Start() {
	p.start = time.Now()

	if !IsInteractive() {
		return
	}

	p.spin = 0

	HideCursor()

	go func() {
		for {
			if p.ctx.Err() != nil {
				return
			}

			p.spin++
			if p.spin > 7 {
				p.spin = 0
			}
			p.print()
			time.Sleep(spinnerDelayMilliseconds * time.Millisecond)
		}
	}()
}

func (p *Progress) Stop() {
	p.cancel()

	if IsInteractive() {
		SavePos()
		fmt.Print(strings.Repeat(" ", p.printedLen+1))
		RestorePos()
		SetColor(color.FgHiBlue)
	}

	fmt.Print(p.elapsed() + " ")
	SetColor(color.Reset)
	ShowCursor()
}
