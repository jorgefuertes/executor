package terminal

import (
	"context"
	"fmt"
	"math"
	"strings"
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
	Action(InfoLevel, p.title, false)
	et := p.elapsed()
	printedLen := len(p.title) + 4 + len(et)
	if p.ctx.Err() == nil {
		color.Set(Blue...)
		fmt.Print(et + " ")
		printedLen++
	} else {
		color.Set(Cyan...)
		fmt.Print(et)
	}
	if p.ctx.Err() == nil {
		color.Set(Yellow...)
		fmt.Print(p.spinner.chars[p.spin])
		printedLen += len(p.spinner.chars[p.spin])

		// remaining spinner characters
		if p.printedLen > printedLen {
			diff := p.printedLen - printedLen
			fmt.Print(strings.Repeat(" ", diff))
			fmt.Print(strings.Repeat("\b", diff))
		}
	}

	color.Unset()
	p.printedLen = printedLen
}

func (p *Progress) Start() {
	Action(InfoLevel, p.title, true)
	time.Sleep(slowPrintDelay)
	fmt.Print("\r")

	p.start = time.Now()

	if !IsInteractive() {
		return
	}

	p.spin = 0

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
