package terminal

import (
	"fmt"
	"time"
)

const (
	Slow           bool   = true
	Fast           bool   = false
	slowPrintDelay        = 10 * time.Millisecond
	cursorChar     string = "█"
)

func Print(colorName colorStyle, slow bool, text string) {
	if !slow || !IsInteractive() {
		if HasColor() {
			print(getColorStyle(colorName).Render(text))
		} else {
			print(text)
		}

		return
	}

	for _, r := range text {
		if HasColor() {
			print(getColorStyle(colorName).Render(string(cursorChar)))
		} else {
			print(string(cursorChar))
		}

		time.Sleep(slowPrintDelay)
		print("\b")
		if HasColor() {
			print(getColorStyle(colorName).Render(string(r)))
		} else {
			print(string(r))
		}
	}
}

func PrintF(colorName colorStyle, slow bool, format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	Print(colorName, slow, text)
}

func caret(level Level) {
	switch level {
	case DebugLevel:
		Print(DebugLevelColor, Fast, ">")
	case InfoLevel:
		Print(InfoLevelColor, Fast, ">")
	case WarnLevel:
		Print(WarnLevelColor, Fast, ">")
	case ErrorLevel:
		Print(ErrorLevelColor, Fast, ">")
	}

	print(" ")
}
