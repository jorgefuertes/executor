package commands

import (
	"fmt"
	"net"

	"executor/internal/config"
	"executor/internal/terminal"
)

func Port(cfg *config.Config) error {
	p := terminal.NewProgress(fmt.Sprintf("%s (%d)", cfg.Desc, cfg.Port), cfg.Style)
	p.Start()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), cfg.Timeout)
	p.Stop(err == nil)
	if err != nil {
		return err
	}
	_ = conn.Close()

	return nil
}
