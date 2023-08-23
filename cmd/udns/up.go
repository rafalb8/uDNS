package main

import (
	"fmt"
	"os"

	"github.com/rafalb8/uDNS/internal"
)

// Creates systemd unit file and starts the service
type UpCmd struct {
	Enable bool `help:"Enable systemd unit on boot" default:"false"`
}

func (u *UpCmd) Run() error {
	// Check if running as root
	if os.Getuid() != 0 {
		return fmt.Errorf("must be run as root")
	}

	internal.InstallService(internal.ConfigPath)
	internal.ServiceControl(internal.StartService)
	if u.Enable {
		internal.ServiceControl(internal.EnableService)
	}
	return nil
}
