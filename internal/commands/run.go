package commands

import (
	"fmt"
	"os"
	"os/exec"

	"executor/internal/config"
	"executor/internal/terminal"
	"github.com/olekukonko/tablewriter"
)

func Run(cfg *config.Config) error {
	terminal.SetNoColor(cfg.NoColor)

	if len(cfg.Command) == 0 {
		return ErrEmptyCommand
	}

	mainEnv, err := getEnv(cfg.EnvFileName, cfg.EnvRecurseLevels)
	if err != nil {
		return err
	}

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

	terminal.Action(terminal.InfoLevel, cfg.Desc)

	progress := terminal.NewProgress(cfg.Style)
	progress.Start()
	o, err := cmd.CombinedOutput()
	progress.Stop()
	terminal.Result(err == nil)

	if err != nil && cfg.ShowAnyOutput() {
		terminal.Line(terminal.WarnLevel, "Failed command: "+cfg.Command)
	}

	if cfg.ShowOutput || (err != nil && cfg.ShowOutputOnError) {
		fmt.Println()
		terminal.Line(terminal.WarnLevel, "Command output:")
		fmt.Printf("OUTPUT: %s\n", string(o))
	}

	return err
}
