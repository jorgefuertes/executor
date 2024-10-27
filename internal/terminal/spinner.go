package terminal

import (
	"strings"
	"time"
)

type spinner struct {
	chars []string
	delay time.Duration
}

var spinners = map[string]spinner{
	"dots":   {chars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧"}, delay: time.Millisecond * 30},
	"arrow":  {chars: []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}, delay: time.Millisecond * 80},
	"star":   {chars: []string{"★", "☆"}, delay: time.Millisecond * 250},
	"circle": {chars: []string{"◐", "◓", "◑", "◒"}, delay: time.Millisecond * 100},
	"square": {chars: []string{"▖", "▘", "▝", "▗"}, delay: time.Millisecond * 100},
	"bullet": {chars: []string{"◎", "◉"}, delay: time.Millisecond * 200},
	"line":   {chars: []string{"⎺", "⎻", "⎼", "⎽", "⎼"}, delay: time.Millisecond * 60},
	"bar":    {chars: []string{`|`, `/`, `-`, `\`}, delay: time.Millisecond * 70},
	"o":      {chars: []string{".", "o", "O", "0", "O", "o", ".", " "}, delay: time.Millisecond * 80},
}

func SpinnerStyles() []string {
	var styles []string
	for k := range spinners {
		styles = append(styles, k)
	}
	return styles
}

func SpinnerStylesString() string {
	return strings.Join(SpinnerStyles(), ", ")
}
