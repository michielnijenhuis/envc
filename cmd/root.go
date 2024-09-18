package cmd

import (
	"github.com/michielnijenhuis/cli"
)

func Execute() {
	app := &cli.Application{
		Name:           "envc",
		Version:        "v2.2.1",
		CatchErrors:    true,
		SingleCommand:  true,
		DefaultCommand: "envc",
		Commands: []*cli.Command{
			Command,
		},
	}

	app.RunExit()
}
