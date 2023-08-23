package main

import "github.com/rafalb8/uDNS/internal"

// Stops the service and/or removes it from systemd
type DownCmd struct {
	Uninstall bool `help:"Uninstall systemd unit file" default:"false"`
}

func (d *DownCmd) Run() error {
	internal.ServiceControl(internal.StopService)
	if d.Uninstall {
		internal.RemoveService()
	}
	internal.CleanResolv()
	return nil
}
