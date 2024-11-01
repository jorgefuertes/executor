package commands

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"executor/internal/config"
	"executor/internal/terminal"

	"github.com/olekukonko/tablewriter"
)

func Run(cfg *config.Config) error {
	terminal.SetNoColor(cfg.NoColor)
	defer terminal.ResetColor()
	defer terminal.ShowCursor()

	if len(cfg.Command) == 0 {
		return ErrEmptyCommand
	}

	mainEnv, _ := getEnv(cfg.EnvFileName, "", cfg.EnvRecurseLevels)

	if cfg.ShowEnv {
		fmt.Println()
		terminal.TableTile("Enviroment")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetCaption(true, fmt.Sprintf("%s: %d vars", cfg.EnvFileName, len(mainEnv)))
		table.SetHeader([]string{"Variable", "Value"})
		for k, v := range mainEnv {
			table.Rich([]string{k, v}, []tablewriter.Colors{{tablewriter.FgCyanColor}, {tablewriter.FgHiYellowColor}})
		}
		table.Render()
		fmt.Println()
	}

	cmd := exec.Command("sh", "-c", cfg.Command)

	cmd.Env = os.Environ()
	for k, v := range mainEnv {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	progress := terminal.NewProgress(cfg.Desc, cfg.Style)
	progress.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT)

	go func() {
		<-ch
		err := cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			terminal.Error(err)
		}
	}()

	o, err := cmd.CombinedOutput()
	progress.Stop(err == nil)

	if err != nil && cfg.ShowAnyOutput() {
		terminal.Line(terminal.WarnLevel, "Failed command: "+cfg.Command)
	}

	if cfg.ShowOutput || (err != nil && cfg.ShowOutputOnError) {
		fmt.Println()
		terminal.Line(terminal.WarnLevel, "Command output:")
		fmt.Print(string(o))
	}

	return err
}
