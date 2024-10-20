package commands

import (
	"bytes"
	"executor/internal/terminal"
	"fmt"
	"os"
	"os/exec"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

func Run(c *cli.Context) error {
	envFileName := c.String("env-file")
	envRecurseLevels := c.Int("env-recurse-level")
	showEnv := c.Bool("show-env")
	showOutput := c.String("show-on-success")
	showOnErr := c.String("show-on-err")
	desc := c.String("desc")
	style := c.String("spinner-style")

	terminal.SetNoColor(c.Bool("no-color"))

	cmdLine := c.String("cmd")
	if len(cmdLine) == 0 {
		return ErrEmptyCommand
	}

	mainEnv, err := getEnv(envFileName, envRecurseLevels)
	if err != nil {
		return err
	}

	if showEnv {
		fmt.Println()
		terminal.TableTile("Enviroment")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetCaption(true, fmt.Sprintf("%s: %d vars", envFileName, len(mainEnv)))
		table.SetHeader([]string{"Variable", "Value"})
		for k, v := range mainEnv {
			table.Rich([]string{k, v}, []tablewriter.Colors{{tablewriter.FgCyanColor}, {tablewriter.FgHiYellowColor}})
		}
		table.Render()
		fmt.Println()
	}

	cmd := exec.Command("sh", "-c", cmdLine)

	cmd.Env = os.Environ()
	for k, v := range mainEnv {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	terminal.Action(terminal.InfoLevel, desc)

	progres := terminal.NewProgress(style)
	progres.Start()
	err = cmd.Run()
	progres.Stop()
	terminal.Result(err == nil)

	if err != nil {
		switch showOnErr {
		case "stdout":
			terminal.ShowOutput(terminal.DebugLevel, "COMMAND STANDARD OUTPUT", stdout)
		case "stderr":
			terminal.ShowOutput(terminal.WarnLevel, "COMMAND ERROR OUTPUT", stderr)
		case "both":
			terminal.ShowOutput(terminal.DebugLevel, "COMMAND STANDARD OUTPUT", stdout)
			terminal.ShowOutput(terminal.WarnLevel, "COMMAND ERROR OUTPUT", stderr)
		}
		return ErrExecutionFailed
	}

	switch showOutput {
	case "stdout":
		terminal.ShowOutput(terminal.DebugLevel, "COMMAND STANDARD OUTPUT", stdout)
	case "stderr":
		terminal.ShowOutput(terminal.WarnLevel, "COMMAND ERROR OUTPUT", stderr)
	case "both":
		terminal.ShowOutput(terminal.DebugLevel, "COMMAND STANDARD OUTPUT", stdout)
		terminal.ShowOutput(terminal.WarnLevel, "COMMAND ERROR OUTPUT", stderr)
	}

	return nil
}
