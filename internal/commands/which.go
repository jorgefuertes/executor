package commands

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jorgefuertes/executor/internal/config"
	"github.com/jorgefuertes/executor/internal/terminal"
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
		if !info.IsDir() && info.Mode().IsRegular() && info.Mode().Perm()&0o100 != 0 {
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
		if !info.IsDir() && info.Mode().IsRegular() && info.Mode().Perm()&0o100 != 0 {
			return true
		}
	}

	return false
}

func Which(cfg *config.Config) error {
	ok := isExecutable(cfg.Command)

	if ok && cfg.Silent {
		return nil
	}

	terminal.SetNoColor(cfg.NoColor)
	desc := "Looking for " + cfg.Command
	terminal.Action(terminal.InfoLevel, desc, true)
	terminal.DashedLine(len(desc) + 4)
	terminal.Result(ok)

	if !ok {
		terminal.Line(terminal.WarnLevel, cfg.NotFoundMsg, false)
		return ErrCommandNotFound
	}

	return nil
}
