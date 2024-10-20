package commands

import (
	"executor/internal/terminal"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func isExecutable(cmd string) bool {
	if cmd == "" {
		return false
	}

	if strings.Contains(cmd, "/") {
		info, err := os.Stat(cmd)
		if err != nil {
			return false
		}
		if !info.IsDir() && info.Mode().IsRegular() && info.Mode().Perm()&0100 != 0 {
			return true
		}
		return false
	}

	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, dir := range paths {
		fullPath := filepath.Join(dir, cmd)
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		if !info.IsDir() && info.Mode().IsRegular() && info.Mode().Perm()&0100 != 0 {
			return true
		}
	}

	return false
}

func Which(c *cli.Context) error {
	cmd := c.String("cmd")

	terminal.Action(terminal.InfoLevel, "Looking for "+cmd)

	ok := isExecutable(cmd)
	terminal.Result(ok)

	if ok {
		return nil
	}

	terminal.Line(terminal.WarnLevel, c.String("not-found-msg"))

	return ErrCommandNotFound
}
