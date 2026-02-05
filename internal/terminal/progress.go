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

func NewProgress(title, style string) *Progress {
	ctx, cancel := context.WithCancel(context.Background())

	s, ok := spinners[style]
	if !ok {
		s = spinners["dots"]
	}

	return &Progress{
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

func (p *Progress) Start() {
	p.lock.Lock()

	Action(InfoLevel, p.title, true)
	time.Sleep(slowPrintDelay)

	p.lock.Unlock()
	p.start = time.Now()

	if !IsInteractive() {
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
			Print(ClockColor, Fast, p.elapsed())
			printedLen := len(p.elapsed()) + 1
			if len(p.bufLen()) > 0 {
				Print(SizeColor, Fast, p.bufLen())
				printedLen += len(p.bufLen())
			}
			print(" ")
			printedLen++

			p.spin++
			if p.spin > len(p.spinner.chars)-1 {
				p.spin = 0
			}

			Print(SpinnerColor, Fast, p.spinner.chars[p.spin])
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

	DashedLine()
	print(strings.Repeat("\b", len(p.elapsed())+5))
	Print(ClockColor, Fast, p.elapsed())
	Result(result)
	ShowCursor()
}
