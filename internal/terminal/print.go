package terminal

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

const Slow bool = true
const Fast bool = false
const slowPrintDelay = 10 * time.Millisecond
const cursorChar string = "█"

var (
	White      = []color.Attribute{color.FgHiWhite}
	Gray       = []color.Attribute{color.FgWhite}
	Blue       = []color.Attribute{color.FgHiBlue}
	Green      = []color.Attribute{color.FgHiGreen}
	Red        = []color.Attribute{color.FgHiRed}
	Yellow     = []color.Attribute{color.FgHiYellow}
	Pink       = []color.Attribute{color.FgHiMagenta}
	Cyan       = []color.Attribute{color.FgCyan}
	RedLabel   = []color.Attribute{color.BgHiRed, color.FgHiWhite}
	GreenLabel = []color.Attribute{color.BgHiGreen, color.FgHiWhite}
	CyanLabel  = []color.Attribute{color.BgCyan, color.FgHiBlack}
)

func Print(cs []color.Attribute, slow bool, text string) {
	color.Set(cs...)
	defer color.Unset()

	if slow && IsInteractive() {
		for _, r := range text {
			fmt.Print(cursorChar)
			time.Sleep(slowPrintDelay)
			fmt.Print("\b" + string(r))
		}

		return
	}

	fmt.Print(text)
}

func PrintF(cs []color.Attribute, slow bool, format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	Print(cs, slow, text)
}

func caret(level Level) {
	switch level {
	case DebugLevel:
		Print(Pink, Fast, ">")
	case InfoLevel:
		Print(Green, Fast, ">")
	case WarnLevel:
		Print(Yellow, Fast, ">")
	case ErrorLevel:
		Print(Red, Fast, ">")
	}

	fmt.Print(" ")
}
