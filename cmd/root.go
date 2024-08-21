package cmd

import (
	"log"

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

	if _, err := app.Run(); err != nil {
		log.Fatalln(err)
	}
}
