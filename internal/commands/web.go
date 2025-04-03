package commands

import (
	"errors"
	"net/http"

	"github.com/jorgefuertes/executor/internal/config"
	"github.com/jorgefuertes/executor/internal/terminal"
)

func Web(cfg *config.Config) error {
	p := terminal.NewProgress(cfg.Desc, cfg.Style)
	p.Start()

	c := http.Client{
		Timeout: cfg.Timeout,
	}

	resp, err := c.Get(cfg.URL)

	p.Stop(err == nil && resp.StatusCode == http.StatusOK)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(resp.StatusCode))
	}

	return nil
}
