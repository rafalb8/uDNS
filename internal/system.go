package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ServiceAction string

const (
	StartService   ServiceAction = "start"
	StopService    ServiceAction = "stop"
	EnableService  ServiceAction = "enable"
	DisableService ServiceAction = "disable"
)

func RemoveService() {
	err := exec.Command("systemctl", "disable", "udns.service").Run()
	if err != nil {
		panic(err)
	}
	err = os.Remove(ServicePath)
	if err != nil {
		panic(err)
	}
	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		panic(err)
	}
}

func InstallService(path string) {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	unit, err := os.OpenFile(ServicePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer unit.Close()

	_, err = unit.WriteString(fmt.Sprintf(ServiceTmpl, executablePath, path))
	if err != nil {
		panic(err)
	}

	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		panic(err)
	}
}

func ServiceControl(action ServiceAction) {
	err := exec.Command("systemctl", string(action), "udns.service").Run()
	if err != nil {
		panic(err)
	}
}

type resLine struct {
	Data    string
	Comment bool
}

func ModifyResolv() {
	// Add 127.0.0.1 to resolv.conf
	resolv, err := os.ReadFile("/etc/resolv.conf")
	if err != nil {
		panic(err)
	}

	lines := []resLine{}

	for _, line := range strings.Split(string(resolv), "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			lines = append(lines, resLine{Data: line, Comment: true})
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 && parts[0] == "nameserver" {
			if parts[1] == "127.0.0.1" {
				continue
			}
			lines = append(lines, resLine{Data: line, Comment: false})
		}
	}

	output, err := os.OpenFile("/etc/resolv.conf", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	i := 0
	for _, line := range lines {
		if line.Comment {
			output.WriteString(line.Data + "\n")
			continue
		}
		i++
		if i == 1 {
			output.WriteString("nameserver 127.0.0.1\n")
		}
		output.WriteString(line.Data + "\n")
	}

	if i < 1 {
		output.WriteString("nameserver 127.0.0.1\n")
	}
}

func CleanResolv() {
	// remove nameserver 127.0.0.1 from resolv.conf
	resolv, err := os.ReadFile("/etc/resolv.conf")
	if err != nil {
		panic(err)
	}

	output, err := os.OpenFile("/etc/resolv.conf", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	for _, line := range strings.Split(string(resolv), "\n") {
		line = strings.TrimSpace(line)
		if line == "nameserver 127.0.0.1" {
			continue
		}
		output.WriteString(line + "\n")
	}
}
