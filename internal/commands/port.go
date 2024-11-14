package commands

import (
	"context"
	"fmt"
	"net"

	"executor/internal/config"
	"executor/internal/terminal"
)

func Port(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	p := terminal.NewProgress(fmt.Sprintf("%s (%d)", cfg.Desc, cfg.Port), cfg.Style)
	p.Start()

	var err error
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), cfg.Timeout)
		if err == nil {
			_ = conn.Close()
			break
		}
		if ctx.Done() != nil {
			break
		}
	}

	p.Stop(err == nil)
	if err != nil {
		return err
	}

	return nil
}
