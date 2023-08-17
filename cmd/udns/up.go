package main

import (
	"fmt"
	"os"

	"github.com/rafalb8/uDNS/internal"
)

// Creates systemd unit file and starts the service
type UpCmd struct {
}

func (u *UpCmd) Run() error {
	// Check if running as root
	if os.Getuid() != 0 {
		return fmt.Errorf("must be run as root")
	}

	internal.InstallService(internal.ConfigPath)
	internal.StartService()
	internal.AppendResolv()
	return nil
}
