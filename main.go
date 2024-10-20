package main

import (
	"executor/internal/commands"
	"executor/internal/terminal"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

type Output string

func (o Output) String() string {
	return string(o)
}

const (
	whichCommandName        = "which"
	runCommandName          = "run"
	ShowOnlyStdout   Output = "stdout"
	ShowOnlyStderr   Output = "stderr"
	ShowBoth         Output = "both"
	ShowNone         Output = "none"
)

var version string

func main() {
	app := &cli.App{
		Name:           "executor",
		Usage:          "Execute commands in fancy way",
		Version:        version,
		DefaultCommand: runCommandName,
		Commands: []*cli.Command{
			{
				Name:  whichCommandName,
				Usage: "Check if a command exists in the system path",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "cmd",
						Aliases: []string{"c"},
						Usage:   "Command to check",
					},
					&cli.StringFlag{
						Name:    "not-found-msg",
						Aliases: []string{"m"},
						Usage:   "Text to show if command not found, typically some install hint",
						Value:   "Command not found, please install it now.",
					},
				},
				Action: commands.Which,
			},
			{
				Name:  runCommandName,
				Usage: "Run a command",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "desc",
						Aliases:  []string{"d"},
						Usage:    "Command description",
						Required: true,
					},
					&cli.BoolFlag{
						Name:    "show-env",
						Aliases: []string{"se"},
						Usage:   "Show enviroment before start",
					},
					&cli.StringFlag{
						Name:    "spinner-style",
						Aliases: []string{"st"},
						Usage:   "Spinner style: dots, arrow, star, circle, square, square-star, line, line-star, bar, o",
						Value:   "dots",
					},
					&cli.StringFlag{
						Name:    "show-on-success",
						Aliases: []string{"os"},
						Usage:   "Show stdout, stderr, both or none",
						Value:   ShowNone.String(),
					},
					&cli.StringFlag{
						Name:    "show-on-err",
						Aliases: []string{"oe"},
						Usage:   "Show stdout, stderr or both",
						Value:   ShowOnlyStderr.String(),
					},
					&cli.StringFlag{
						Name:    "env-file",
						Aliases: []string{"n"},
						Usage:   "Enviroment file ('none' to disable)",
						Value:   ".env",
					},
					&cli.IntFlag{
						Name: "env-recurse-levels",
						Aliases: []string{
							"r",
						},
						Usage: "How many levels we should recurse back looking for the env file, if its not an absolute path",
						Value: 5,
					},
					&cli.BoolFlag{
						Name:       "no-color",
						Aliases:    []string{"nc"},
						Value:      false,
						Usage:      "Disable color output and spinner",
						HasBeenSet: true,
					},
					&cli.StringFlag{
						Name:     "cmd",
						Aliases:  []string{"c"},
						Usage:    "Command to run",
						Required: true,
					},
				},
				Action: commands.Run,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println()
		terminal.Error(err)
		fmt.Println()
		os.Exit(1)
	}
}
