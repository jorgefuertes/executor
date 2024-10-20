package terminal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

type SpinnerStyle string

type Progress struct {
	spin       int
	start      time.Time
	elapsed    time.Duration
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
		elapsed:    time.Duration(0),
		ctx:        ctx,
		cancel:     cancel,
		chars:      chars,
		printedLen: 0,
	}
}

func (p *Progress) print() {
	SavePos()
	p.elapsed = time.Since(p.start)
	spinStr := fmt.Sprintf("[Elapsed %.0fm%.0fs] %s", p.elapsed.Minutes(), p.elapsed.Seconds(), p.chars[p.spin])
	fmt.Print(spinStr)
	p.printedLen = len(spinStr)
	RestorePos()
}

func (p *Progress) Start() {
	p.spin = 0
	p.start = time.Now()
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
			time.Sleep(150 * time.Millisecond)
		}
	}()
}

func (p *Progress) Stop() {
	p.cancel()

	SavePos()
	fmt.Print(strings.Repeat(" ", p.printedLen+1))
	RestorePos()
	color.Set(color.FgHiBlue)
	fmt.Printf("[Total %.0fm%.0fs] ", p.elapsed.Minutes(), p.elapsed.Seconds())
	color.Set(color.Reset)
	ShowCursor()
}
