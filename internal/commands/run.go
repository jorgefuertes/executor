package commands

import (
	"bytes"
	"executor/internal/terminal"
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

func Run(c *cli.Context) error {
	envFile := c.String("env-file")
	showEnv := c.Bool("show-env")
	showOutput := c.String("show-on-success")
	showOnErr := c.String("show-on-err")
	desc := c.String("desc")

	cmdLine := c.String("cmd")
	if len(cmdLine) == 0 {
		return ErrEmptyCommand
	}

	mainEnv, err := godotenv.Read(envFile)
	if err != nil {
		curDir, _ := os.Getwd()
		terminal.Error(fmt.Errorf("cannot read %s/%s", curDir, envFile))
		terminal.Error(err)
		os.Exit(1)
	}

	if showEnv {
		fmt.Println()
		terminal.TableTile("Enviroment")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetCaption(true, fmt.Sprintf("%s: %d vars", envFile, len(mainEnv)))
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

	progres := terminal.NewProgress()
	progres.Start()

	err = cmd.Start()
	if err != nil {
		terminal.Error(err)
		os.Exit(1)
	}

	err = cmd.Wait()
	progres.Stop()

	if err == nil {
		switch showOutput {
		case "stdout":
			terminal.ShowOutput(terminal.DebugLevel, "COMMAND STANDARD OUTPUT", stdout)
		case "stderr":
			terminal.ShowOutput(terminal.WarnLevel, "COMMAND ERROR OUTPUT", stderr)
		case "both":
			terminal.ShowOutput(terminal.DebugLevel, "COMMAND STANDARD OUTPUT", stdout)
			terminal.ShowOutput(terminal.WarnLevel, "COMMAND ERROR OUTPUT", stderr)
		}
	} else {
		switch showOnErr {
		case "stdout":
			terminal.ShowOutput(terminal.DebugLevel, "COMMAND STANDARD OUTPUT", stdout)
		case "stderr":
			terminal.ShowOutput(terminal.WarnLevel, "COMMAND ERROR OUTPUT", stderr)
		case "both":
			terminal.ShowOutput(terminal.DebugLevel, "COMMAND STANDARD OUTPUT", stdout)
			terminal.ShowOutput(terminal.WarnLevel, "COMMAND ERROR OUTPUT", stderr)
		}
	}

	return ErrExecutionFailed
}
