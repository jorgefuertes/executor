//go:build linux || darwin

package terminal

import (
	"os"

	"golang.org/x/sys/unix"
)

func (t *Term) DisableEcho() {
	if !t.interactive {
		return
	}

	fd := int(os.Stdin.Fd())

	termios, err := unix.IoctlGetTermios(fd, ioctlGetTermios)
	if err != nil {
		return
	}

	saved := *termios
	t.savedTermios = &saved

	termios.Lflag &^= unix.ECHO
	_ = unix.IoctlSetTermios(fd, ioctlSetTermios, termios)
}

func (t *Term) RestoreEcho() {
	termios, ok := t.savedTermios.(*unix.Termios)
	if !ok {
		return
	}

	fd := int(os.Stdin.Fd())
	_ = unix.IoctlSetTermios(fd, ioctlSetTermios, termios)
	t.savedTermios = nil
}
