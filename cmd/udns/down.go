package main

// Stops the service and/or removes it from systemd
type DownCmd struct {
	Uninstall bool `arg:"" help:"Uninstall systemd unit file"`
}

func (d *DownCmd) Run() error {
	return nil
}
