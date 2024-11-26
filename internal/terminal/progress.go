package terminal

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
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
	lock       *sync.Mutex
	OutBuffer  *bytes.Buffer
	ErrBuffer  *bytes.Buffer
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
		lock:       &sync.Mutex{},
		OutBuffer:  new(bytes.Buffer),
		ErrBuffer:  new(bytes.Buffer),
	}
}

func (p *Progress) elapsed() string {
	t := time.Since(p.start)
	m := int(math.Floor(t.Minutes()))
	s := int(math.Floor(t.Seconds())) % 60
	ms := t.Milliseconds()

	out := "["
	if m > 0 {
		out += fmt.Sprintf("%d min", m)
	}
	if s > 0 {
		if m > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%d sec", s)
	}
	if m == 0 && s == 0 {
		out += fmt.Sprintf("%d ms", ms)
	}
	out += "]"

	return out
}

func (p *Progress) bufLen() string {
	if p.OutBuffer.Len() > 0 || p.ErrBuffer.Len() > 0 {
		return fmt.Sprintf("[%s]", humanize.Bytes(uint64(p.OutBuffer.Len()+p.ErrBuffer.Len())))
	}

	return ""
}

func (p *Progress) print() {
	p.lock.Lock()
	defer p.lock.Unlock()
	fmt.Print("\r")
	Action(InfoLevel, p.title, false)

	et := p.elapsed()
	bl := p.bufLen()
	printedLen := len(p.title) + 4 + len(et)

	if p.ctx.Err() == nil {
		color.Set(Blue...)
		fmt.Print(et)
		printedLen++
		if len(bl) > 0 {
			printedLen += 1 + len(bl)
			color.Set(DarkGreen...)
			fmt.Print(bl)
		}
		fmt.Print(" ")
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
	p.lock.Lock()
	Action(InfoLevel, p.title, true)
	time.Sleep(slowPrintDelay)
	fmt.Print("\r")
	p.lock.Unlock()

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
