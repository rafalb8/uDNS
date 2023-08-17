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
