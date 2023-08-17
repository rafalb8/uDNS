package main

import (
	"github.com/rafalb8/uDNS/internal/server"
)

// Start the server
type ServerCmd struct {
	Port string `arg:"" help:"DNS server port" default:"53"`
	HttpPort string `arg:"" help:"HTTP server port" default:"8367"`
}

func (s *ServerCmd) Run() error {
	return server.Start(s.Port, s.HttpPort)
}
