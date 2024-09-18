package cmd

import (
	"github.com/michielnijenhuis/cli"
)

func Execute() {
	app := &cli.Application{
		Name:           "envc",
		Version:        "v2.2.3",
		CatchErrors:    true,
		SingleCommand:  true,
		DefaultCommand: "envc",
		Commands: []*cli.Command{
			Command,
		},
	}

	app.RunExit()
}
