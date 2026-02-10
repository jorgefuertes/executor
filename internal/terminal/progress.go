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
)

type Progress struct {
	t              *Term
	title          string
	spinner        spinner
	lastPrintedLen int
	spin           int
	start          time.Time
	ctx            context.Context
	cancel         context.CancelFunc
	lock           *sync.Mutex
	OutBuffer      *bytes.Buffer
	ErrBuffer      *bytes.Buffer
}

func (t *Term) NewProgress(title string) *Progress {
	ctx, cancel := context.WithCancel(context.Background())

	s, ok := spinners[t.cfg.Style]
	if !ok {
		s = spinners["dots"]
	}

	return &Progress{
		t:         t,
		title:     title,
		spinner:   s,
		start:     time.Now(),
		ctx:       ctx,
		cancel:    cancel,
		lock:      &sync.Mutex{},
		OutBuffer: new(bytes.Buffer),
		ErrBuffer: new(bytes.Buffer),
	}
}

func (p *Progress) elapsed() string {
	et := time.Since(p.start)
	m := int(math.Floor(et.Minutes()))
	s := int(math.Floor(et.Seconds())) % 60
	ms := et.Milliseconds()

	var b strings.Builder

	b.WriteString("[")
	if m > 0 {
		b.WriteString(fmt.Sprintf("%d min", m))
	}
	if s > 0 {
		if m > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("%d sec", s))
	}
	if m == 0 && s == 0 {
		b.WriteString(fmt.Sprintf("%d ms", ms))
	}
	b.WriteString("]")

	return b.String()
}

func (p *Progress) bufLen() string {
	if p.OutBuffer.Len() > 0 || p.ErrBuffer.Len() > 0 {
		return fmt.Sprintf("[%s]", humanize.Bytes(uint64(p.OutBuffer.Len()+p.ErrBuffer.Len())))
	}

	return ""
}

func (p *Progress) Start() {
	p.lock.Lock()

	p.t.Action(InfoLevel, p.title, true)
	time.Sleep(slowPrintDelay)

	p.lock.Unlock()
	p.start = time.Now()

	if !p.t.IsInteractive() {
		return
	}

	p.spin = 0
	go func() {
		for {
			p.lock.Lock()

			if p.ctx.Err() != nil {
				p.lock.Unlock()
				return
			}

			print(" ")
			p.t.Print(ClockColor, Fast, p.elapsed())
			printedLen := len(p.elapsed()) + 1
			if len(p.bufLen()) > 0 {
				p.t.Print(SizeColor, Fast, p.bufLen())
				printedLen += len(p.bufLen())
			}
			print(" ")
			printedLen++

			p.spin++
			if p.spin > len(p.spinner.chars)-1 {
				p.spin = 0
			}

			p.t.Print(SpinnerColor, Fast, p.spinner.chars[p.spin])
			printedLen++
			if p.lastPrintedLen > printedLen {
				print(strings.Repeat(" ", p.lastPrintedLen-printedLen))
				print(strings.Repeat("\b", p.lastPrintedLen-printedLen))
			}
			p.lastPrintedLen = printedLen

			print(strings.Repeat("\b", printedLen))

			p.lock.Unlock()
			time.Sleep(p.spinner.delay)
		}
	}()
}

func (p *Progress) Stop(result bool) {
	p.cancel()
	defer p.t.ShowCursor()

	if !p.t.IsInteractive() {
		p.t.Print(SecondaryColor, Fast, ellipsis)
		p.t.Print(ClockColor, Fast, p.elapsed())
		p.t.Print(SecondaryColor, Fast, ellipsis)
		p.t.Result(result)

		return
	}

	p.t.DashedLine()
	print(strings.Repeat("\b", len(p.elapsed())+5))
	p.t.Print(ClockColor, Fast, p.elapsed())
	p.t.Result(result)
}
