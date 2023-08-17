package main

import (
	"github.com/alecthomas/kong"
)

var cli struct {
	Version VersionCmd `cmd:"" help:"Show version"`
	Edit    EditCmd    `cmd:"" help:"Open VS Code inside config directory"`
	Up      UpCmd      `cmd:"" help:"Start µDNS"`
	Down    DownCmd    `cmd:"" help:"Stop µDNS"`
	Server  ServerCmd  `cmd:"" help:"Run µDNS server"`
	Test    TestCmd    `cmd:"" help:"Test resolver"`
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
