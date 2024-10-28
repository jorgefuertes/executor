package terminal

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/fatih/color"
)

type Progress struct {
	title      string
	spinner    spinner
	spin       int
	start      time.Time
	printedLen int
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewProgress(title, style string) *Progress {
	ctx, cancel := context.WithCancel(context.Background())

	s, ok := spinners[style]
	if !ok {
		s = spinners["dots"]
	}

	return &Progress{
		title:      title,
		spinner:    s,
		spin:       0,
		start:      time.Now(),
		printedLen: 0,
		ctx:        ctx,
		cancel:     cancel,
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
	fmt.Print("\r")
	Action(InfoLevel, p.title)
	et := p.elapsed()
	p.printedLen = len(p.title) + 4 + len(et)
	if p.ctx.Err() == nil {
		SetColor(color.FgHiBlue)
		fmt.Print(et + " ")
		p.printedLen++
	} else {
		SetColor(color.FgCyan)
		fmt.Print(et)
	}
	if p.ctx.Err() == nil {
		SetColor(color.FgYellow)
		fmt.Print(p.spinner.chars[p.spin])
		p.printedLen += len(p.spinner.chars[p.spin])
	}

	ResetColor()
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

func (p *Progress) Stop(ok bool) {
	p.cancel()
	p.print()
	if IsInteractive() {
		DashedLine(p.printedLen)
	}

	Result(ok)
	ShowCursor()
}
