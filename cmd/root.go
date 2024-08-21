package cmd

import (
	"github.com/michielnijenhuis/cli"
)

func Execute() {
	app := &cli.Application{
		Name:           "envc",
		Version:        "v1.0.0",
		CatchErrors:    true,
		AutoExit:       true,
		SingleCommand:  true,
		DefaultCommand: "envc",
	}

	app.Add(Command)

	app.Run()
}