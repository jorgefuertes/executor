package commands

import (
	"context"
	"fmt"
	"net"

	"github.com/jorgefuertes/executor/internal/config"
	"github.com/jorgefuertes/executor/internal/terminal"
)

func Port(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	t := terminal.New(cfg)
	defer t.CleanUp()

	p := t.NewProgress(fmt.Sprintf("%s (%d)", cfg.Desc, cfg.Port))
	p.Start()

	var err error
	var conn net.Conn
	for {
		conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), cfg.Timeout)
		if err == nil {
			_ = conn.Close()
			break
		}
		if ctx.Err() != nil {
			err = fmt.Errorf("timeout after %v", cfg.Timeout)
			break
		}
	}

	p.Stop(err == nil)
	if err != nil {
		t.Error(err)
		return err
	}

	return nil
}
