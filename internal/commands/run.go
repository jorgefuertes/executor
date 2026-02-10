package commands

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/jorgefuertes/executor/internal/config"
	"github.com/jorgefuertes/executor/internal/terminal"
	"github.com/olekukonko/tablewriter"
)

func Run(cfg *config.Config) error {
	t := terminal.New(cfg)
	defer t.CleanUp()

	if len(cfg.Command) == 0 {
		return ErrEmptyCommand
	}

	mainEnv, _ := getEnv(cfg.EnvFileName, "", cfg.EnvRecurseLevels)

	if cfg.ShowEnv {
		fmt.Println()
		t.TableTile("Environment")
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

	progress := t.NewProgress(cfg.Desc)
	cmd.Stdout = progress.OutBuffer
	cmd.Stderr = progress.ErrBuffer
	progress.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT)

	go func() {
		<-ch
		err := cmd.Process.Signal(syscall.SIGTERM)
		progress.Stop(false)
		if err != nil {
			t.Error(err)
		}
		os.Exit(1)
	}()

	err := cmd.Run()
	progress.Stop(err == nil)

	if err != nil && cfg.ShowAnyOutput() {
		t.Line(terminal.WarnLevel, "Failed command: "+cfg.Command, false)
	}

	if cfg.ShowOutput || (err != nil && cfg.ShowOutputOnError) {
		fmt.Println()
		t.Line(terminal.WarnLevel, "Command output:", false)
		fmt.Print(progress.OutBuffer.String())
		fmt.Println()
		fmt.Print(progress.ErrBuffer.String())
	}

	return err
}
