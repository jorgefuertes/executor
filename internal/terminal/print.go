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

func (t *Term) Print(colorName colorStyle, slow bool, text string) {
	if !slow || !t.IsInteractive() {
		if t.HasColor() {
			print(getColorStyle(colorName).Render(text))
		} else {
			print(text)
		}

		return
	}

	for _, r := range text {
		if t.HasColor() {
			print(getColorStyle(colorName).Render(string(cursorChar)))
		} else {
			print(string(cursorChar))
		}

		time.Sleep(slowPrintDelay)
		print("\b")
		if t.HasColor() {
			print(getColorStyle(colorName).Render(string(r)))
		} else {
			print(string(r))
		}
	}
}

func (t *Term) PrintF(colorName colorStyle, slow bool, format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	t.Print(colorName, slow, text)
}

func (t *Term) caret(level Level) {
	switch level {
	case DebugLevel:
		t.Print(DebugLevelColor, Fast, ">")
	case InfoLevel:
		t.Print(InfoLevelColor, Fast, ">")
	case WarnLevel:
		t.Print(WarnLevelColor, Fast, ">")
	case ErrorLevel:
		t.Print(ErrorLevelColor, Fast, ">")
	}

	print(" ")
}

func (t *Term) Line(level Level, msg string, slow bool) {
	t.caret(level)
	t.Print(PrimaryColor, slow, msg+"\n")
}

func (t *Term) Action(level Level, msg string, slow bool) int {
	t.caret(level)

	t.Print(PrimaryColor, slow, msg+":")

	return len(msg) + 3
}

func (t *Term) Error(err error) {
	if err == nil {
		return
	}

	t.caret(ErrorLevel)
	t.Print(ErrorColor, false, "ERROR")
	t.PrintF(PrimaryColor, false, ": %s", err.Error())
	fmt.Println()
}
