package main

import (
	"os"
	"os/exec"

	"github.com/rafalb8/uDNS/internal"
)

type EditCmd struct {
}

func (e *EditCmd) Run() error {
	cmd := exec.Command("code", internal.ConfigPath)
	cmd.Env = os.Environ()
	return cmd.Run()
}
