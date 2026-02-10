package config

import (
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

type Config struct {
	EnvFileName       string
	EnvRecurseLevels  int
	ShowEnv           bool
	Style             string
	NoColor           bool
	NoInteractive     bool
	ShowOutput        bool
	ShowOutputOnError bool
	Silent            bool
	NotFoundMsg       string
	Desc              string
	Command           string
	Host              string
	Port              int
	URL               string
	Timeout           time.Duration
}

func New(c *cli.Context) *Config {
	cfg := &Config{
		EnvFileName:       c.String("env-file"),
		EnvRecurseLevels:  c.Int("env-recurse-levels"),
		ShowEnv:           c.Bool("show-env"),
		Style:             c.String("spinner-style"),
		NoColor:           c.Bool("no-color"),
		NoInteractive:     c.Bool("no-interactive"),
		ShowOutput:        c.Bool("show-output"),
		ShowOutputOnError: c.Bool("show-output-on-error"),
		Silent:            c.Bool("silent"),
		NotFoundMsg:       c.String("not-found-msg"),
		Desc:              c.String("desc"),
		Command:           c.String("cmd"),
		Host:              c.String("host"),
		Port:              c.Int("port"),
		URL:               c.String("url"),
		Timeout:           time.Duration(c.Int("timeout")) * time.Second,
	}

	if cfg.EnvRecurseLevels < 1 {
		cfg.EnvRecurseLevels = 5
	}

	if cfg.Timeout < 1 {
		cfg.Timeout = 5
	}

	return cfg
}

func (c Config) Print() {
	print("\n*** Configuration ***\n\n")

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Flag", "Value"})
	t.SetColumnColor(
		tablewriter.Colors{tablewriter.FgCyanColor, tablewriter.Italic},
		tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
	)

	t.AppendBulk([][]string{
		{"EnvFileName", c.EnvFileName},
		{"EnvRecurseLevels", fmt.Sprintf("%d", c.EnvRecurseLevels)},
		{"ShowEnv", fmt.Sprintf("%t", c.ShowEnv)},
		{"Style", c.Style},
		{"NoColor", fmt.Sprintf("%t", c.NoColor)},
		{"NoInteractive", fmt.Sprintf("%t", c.NoInteractive)},
		{"ShowOutput", fmt.Sprintf("%t", c.ShowOutput)},
		{"ShowOutputOnError", fmt.Sprintf("%t", c.ShowOutputOnError)},
		{"Silent", fmt.Sprintf("%t", c.Silent)},
		{"NotFoundMsg", c.NotFoundMsg},
		{"Desc", c.Desc},
		{"Command", c.Command},
		{"Host", c.Host},
		{"Port", fmt.Sprintf("%d", c.Port)},
		{"URL", c.URL},
		{"Timeout", fmt.Sprintf("%d", c.Timeout)},
	})

	t.Render()
}

func (c Config) ShowAnyOutput() bool {
	return c.ShowOutput || c.ShowOutputOnError
}
