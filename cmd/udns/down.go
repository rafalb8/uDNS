package main

// Stops the service and/or removes it from systemd
type DownCmd struct {
	Uninstall bool `help:"Uninstall systemd unit file" default:"false"`
}

func (d *DownCmd) Run() error {
	return nil
}
