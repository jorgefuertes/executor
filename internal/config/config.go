package config

import (
	"fmt"
	"os"

	"executor/internal/terminal"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

type Config struct {
	EnvFileName       string
	EnvRecurseLevels  int
	ShowEnv           bool
	Style             string
	NoColor           bool
	ShowOutput        bool
	ShowOutputOnError bool
	Silent            bool
	NotFoundMsg       string
	Desc              string
	Command           string
}

func New(c *cli.Context) *Config {
	return &Config{
		EnvFileName:       c.String("env-file"),
		EnvRecurseLevels:  c.Int("env-recurse-level"),
		ShowEnv:           c.Bool("show-env"),
		Style:             c.String("spinner-style"),
		NoColor:           c.Bool("no-color"),
		ShowOutput:        c.Bool("show-output"),
		ShowOutputOnError: c.Bool("show-output-on-error"),
		Silent:            c.Bool("silent"),
		NotFoundMsg:       c.String("not-found-msg"),
		Desc:              c.String("desc"),
		Command:           c.String("cmd"),
	}
}

func (c Config) Print() {
	terminal.TableTile("Config")
	t := tablewriter.NewWriter(os.Stdout)
	t.AppendBulk([][]string{
		{"EnvFileName", c.EnvFileName},
		{"EnvRecurseLevels", fmt.Sprintf("%d", c.EnvRecurseLevels)},
		{"ShowEnv", fmt.Sprintf("%t", c.ShowEnv)},
		{"Style", c.Style},
		{"NoColor", fmt.Sprintf("%t", c.NoColor)},
		{"ShowOutput", fmt.Sprintf("%t", c.ShowOutput)},
		{"ShowOutputOnError", fmt.Sprintf("%t", c.ShowOutputOnError)},
		{"Silent", fmt.Sprintf("%t", c.Silent)},
		{"NotFoundMsg", c.NotFoundMsg},
		{"Desc", c.Desc},
		{"Command", c.Command},
	})
	t.Render()
}

func (c Config) ShowAnyOutput() bool {
	return c.ShowOutput || c.ShowOutputOnError
}
