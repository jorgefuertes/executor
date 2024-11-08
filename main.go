package main

import (
	"fmt"
	"os"

	"executor/internal/commands"
	"executor/internal/config"
	"executor/internal/terminal"

	"github.com/urfave/cli/v2"
)

type Output string

func (o Output) String() string {
	return string(o)
}

const (
	whichCommandName = "which"
	runCommandName   = "run"
	portCommandName  = "port"
	webCommandName   = "web"
)

var version string

func main() {
	app := &cli.App{
		Name:           "executor",
		Usage:          "Execute commands in fancy way",
		Version:        version,
		DefaultCommand: runCommandName,
		CommandNotFound: func(c *cli.Context, command string) {
			terminal.Error(fmt.Errorf("command not found: %s", command))
			terminal.CleanUp()
			os.Exit(1)
		},
		Commands: []*cli.Command{
			{
				Name:  whichCommandName,
				Usage: "Check if a command exists in the system path",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "silent",
						Aliases: []string{"s"},
						Usage:   "Silent if command is found",
						Value:   false,
					},
					&cli.BoolFlag{
						Name:       "no-color",
						Aliases:    []string{"nc"},
						Value:      false,
						Usage:      "Disable color output and spinner",
						HasBeenSet: true,
					},
					&cli.BoolFlag{
						Name:    "show-config",
						Aliases: []string{"sc"},
						Usage:   "Show config before start",
						Value:   false,
					},
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
				Action: newActionFunc(commands.Which),
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
						Value:   false,
					},
					&cli.BoolFlag{
						Name:    "show-config",
						Aliases: []string{"sc"},
						Usage:   "Show config before start",
					},
					&cli.StringFlag{
						Name:    "spinner-style",
						Aliases: []string{"st"},
						Usage:   "Spinner style: " + terminal.SpinnerStylesString(),
						Value:   "dots",
					},
					&cli.BoolFlag{
						Name:    "show-output",
						Aliases: []string{"so"},
						Usage:   "Show command output when command it's successful",
						Value:   false,
					},
					&cli.BoolFlag{
						Name:    "show-output-on-error",
						Aliases: []string{"soe"},
						Usage:   "Set false to hide command output when command it's not successful",
						Value:   true,
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
				Action: newActionFunc(commands.Run),
			},
			{
				Name:  portCommandName,
				Usage: "Check if a port is open",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "desc",
						Aliases:  []string{"d"},
						Usage:    "Port check description",
						Required: true,
					},
					&cli.StringFlag{Name: "host", Aliases: []string{"i"}, Usage: "Host to check", Value: "localhost"},
					&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Usage: "Port to check", Required: true},
					&cli.IntFlag{Name: "timeout", Aliases: []string{"t"}, Usage: "Timeout in seconds", Value: 5},
				},
				Action: newActionFunc(commands.Port),
			},
			{
				Name:  webCommandName,
				Usage: "Check if a web page is running and responding successfully",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "desc",
						Aliases:  []string{"d"},
						Usage:    "URL check description",
						Required: true,
					},
					&cli.StringFlag{Name: "url", Aliases: []string{"u"}, Usage: "URL to check", Required: true},
					&cli.IntFlag{Name: "timeout", Aliases: []string{"t"}, Usage: "Timeout in seconds", Value: 5},
				},
				Action: newActionFunc(commands.Web),
			},
		},
	}

	defer terminal.CleanUp()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println()
		terminal.Error(err)
		fmt.Println()
		terminal.CleanUp()
		os.Exit(1)
	}
}

type actionFunc func(cfg *config.Config) error

func newActionFunc(fn actionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		cfg := config.New(c)

		if c.Bool("show-config") {
			cfg.Print()
		}

		return fn(cfg)
	}
}
