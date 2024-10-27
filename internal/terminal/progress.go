package terminal

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/color"
)

const spinnerDelayMilliseconds = 50

type Progress struct {
	spinner    spinner
	spin       int
	start      time.Time
	ctx        context.Context
	cancel     context.CancelFunc
	printedLen int
}

func NewProgress(style string) *Progress {
	ctx, cancel := context.WithCancel(context.Background())

	s, ok := spinners[style]
	if !ok {
		s = spinners["dots"]
	}

	return &Progress{
		spinner:    s,
		spin:       0,
		start:      time.Now(),
		ctx:        ctx,
		cancel:     cancel,
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
	p.printedLen = len(et + " " + p.spinner.chars[p.spin])
	SetColor(color.FgCyan)
	fmt.Print(et + " ")
	SetColor(color.FgYellow)
	fmt.Print(p.spinner.chars[p.spin])
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
			if p.spin > len(p.spinner.chars)-1 {
				p.spin = 0
			}
			p.print()
			time.Sleep(p.spinner.delay)
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
