//go:build windows

package terminal

// Windows does not use termios; echo control is a no-op.
func (t *Term) DisableEcho() {}

func (t *Term) RestoreEcho() {}
