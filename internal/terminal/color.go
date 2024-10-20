package terminal

import "github.com/fatih/color"

func SetColor(p ...color.Attribute) {
	if IsInteractive() {
		color.Set(p...)
	}
}

func ResetColor() {
	SetColor(color.Reset)
}
