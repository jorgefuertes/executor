package config

import (
	"time"

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
	Host              string
	Port              int
	URL               string
	Timeout           time.Duration
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
		Host:              c.String("host"),
		Port:              c.Int("port"),
		URL:               c.String("url"),
		Timeout:           time.Duration(c.Int("timeout")) * time.Second,
	}
}

func (c Config) ShowAnyOutput() bool {
	return c.ShowOutput || c.ShowOutputOnError
}
