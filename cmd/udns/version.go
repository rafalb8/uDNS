package main

import (
	"fmt"
	"runtime/debug"
)

var Version string

type VersionCmd struct {
}

func (v *VersionCmd) Run() error {
	settings := map[string]string{}
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			settings[setting.Key] = setting.Value
		}
	}

	// Show version
	fmt.Println("rEnv version:", Version)
	fmt.Println("Commit time:", settings["vcs.time"])
	fmt.Println("Commit hash:", settings["vcs.revision"])

	return nil
}
