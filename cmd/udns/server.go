package main

import (
	"github.com/rafalb8/uDNS/internal/server"
)

// Start the server
type ServerCmd struct {
	Config   string `help:"Path to config dir"`
	NoResolv bool   `help:"Disable resolv.conf modification" default:"false"`

	Port     string `arg:"" help:"DNS server port" default:"53"`
	HttpPort string `arg:"" help:"HTTP server port" default:"8367"`
}

func (s *ServerCmd) Run() error {
	return server.Start(s.Port, s.HttpPort, s.Config, s.NoResolv)
}
