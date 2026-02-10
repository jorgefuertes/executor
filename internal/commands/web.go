package commands

import (
	"errors"
	"net/http"

	"github.com/jorgefuertes/executor/internal/config"
	"github.com/jorgefuertes/executor/internal/terminal"
)

func Web(cfg *config.Config) error {
	t := terminal.New(cfg)
	defer t.CleanUp()

	p := t.NewProgress(cfg.Desc)
	p.Start()

	c := http.Client{
		Timeout: cfg.Timeout,
	}

	resp, err := c.Get(cfg.URL)

	p.Stop(err == nil && resp.StatusCode == http.StatusOK)

	if err != nil {
		t.Error(err)

		return err
	}

	if resp.StatusCode != http.StatusOK {
		err := errors.New(http.StatusText(resp.StatusCode))
		t.Error(err)

		return err
	}

	return nil
}
