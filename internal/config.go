package internal

import (
	"os"
	"path"
	"strings"
)

var ConfigPath = func() string {
	user := os.Getenv("SUDO_USER")
	if user == "" {
		user = os.Getenv("USER")
	}

	file, err := os.ReadFile("/etc/passwd")
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(string(file), "\n") {
		parts := strings.Split(line, ":")
		if len(parts) > 0 && parts[0] == user {
			return path.Join(parts[5], ".config", "udns")
		}
	}
	panic("User Home not found")
}()

const ServicePath = "/usr/lib/systemd/system/udns.service"
const ServiceTmpl = `
[Unit]
Description=uDNS service
After=network.target

[Service]
ExecStart=%s server --config %s
Restart=on-failure
User=root

[Install]
WantedBy=multi-user.target`
