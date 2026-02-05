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
	ellipsis              = "…"
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

func Line(level Level, msg string, slow bool) {
	caret(level)
	Print(PrimaryColor, slow, msg+"\n")
}

func Action(level Level, msg string, slow bool) int {
	caret(level)

	Print(PrimaryColor, slow, msg+":")

	return len(msg) + 3
}

func Error(err error) {
	if err == nil {
		return
	}

	caret(ErrorLevel)
	Print(ErrorColor, false, "ERROR")
	PrintF(PrimaryColor, false, ": %s", err.Error())
	fmt.Println()
}
