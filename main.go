package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"executor/internal/terminal"
	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
)

type Output string

func (o Output) String() string {
	return string(o)
}

const (
	ShowOnlyStdout Output = "stdout"
	ShowOnlyStderr Output = "stderr"
	ShowBoth       Output = "both"
	ShowNone       Output = "none"
)

var version string

func main() {
	// flags
	desc := flag.String("desc", "", "Command description")
	showEnv := flag.Bool("show-env", false, "Show enviroment before start")
	showOutput := flag.String("show-output", ShowNone.String(), "Show stdout, stderr, both or none")
	showOnErr := flag.String("show-on-err", ShowOnlyStderr.String(), "Show stdout, stderr or both")
	envFile := flag.String("env-file", "./.env.properties", "Enviroment file")
	showVersion := flag.Bool("version", false, "Show version")

	flag.Parse()

	usage := func() {
		flag.Usage()
		fmt.Printf("\n[*] Include command and arguments at the end of the command line\n\n")
		os.Exit(1)
	}

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *desc == "" {
		usage()
	}

	args := flag.Args()
	if len(args) == 0 {
		usage()
	}

	cmdLine := strings.Join(args, " ")
	cmdLine = strings.Trim(cmdLine, `"`)

	if cmdLine == "" {
		usage()
	}

	mainEnv, err := godotenv.Read(*envFile)
	if err != nil {
		curDir, _ := os.Getwd()
		terminal.Error(fmt.Errorf("cannot read %s/%s", curDir, *envFile))
		terminal.Error(err)
		os.Exit(1)
	}

	if *showEnv {
		fmt.Println()
		terminal.TableTile("Enviroment")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetCaption(true, fmt.Sprintf("%s: %d vars", *envFile, len(mainEnv)))
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

	terminal.Action(terminal.InfoLevel, *desc)

	progres := terminal.NewProgress()
	progres.Start()

	err = cmd.Start()
	if err != nil {
		terminal.Error(err)
		os.Exit(1)
	}

	err = cmd.Wait()
	progres.Stop()

	terminal.ActionError(err)
	if err == nil {
		switch Output(*showOutput) {
		case ShowOnlyStdout:
			terminal.ShowOutput(terminal.DebugLevel, "COMMAND STANDARD OUTPUT", stdout)
		case ShowOnlyStderr:
			terminal.ShowOutput(terminal.WarnLevel, "COMMAND ERROR OUTPUT", stderr)
		case ShowBoth:
			terminal.ShowOutput(terminal.DebugLevel, "COMMAND STANDARD OUTPUT", stdout)
			terminal.ShowOutput(terminal.WarnLevel, "COMMAND ERROR OUTPUT", stderr)
		}
	} else {
		switch Output(*showOnErr) {
		case ShowOnlyStdout:
			terminal.ShowOutput(terminal.DebugLevel, "COMMAND STANDARD OUTPUT", stdout)
		case ShowOnlyStderr:
			terminal.ShowOutput(terminal.WarnLevel, "COMMAND ERROR OUTPUT", stderr)
		case ShowBoth:
			terminal.ShowOutput(terminal.DebugLevel, "COMMAND STANDARD OUTPUT", stdout)
			terminal.ShowOutput(terminal.WarnLevel, "COMMAND ERROR OUTPUT", stderr)
		}
	}

	os.Exit(cmd.ProcessState.ExitCode())
}
