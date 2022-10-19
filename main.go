package main

import (
	"log"
	"os"

	"github.com/labbs/alfred/cmd"
	"github.com/urfave/cli/v2"
)

var (
	version = "development"
)

func main() {
	app := cli.NewApp()
	app.Name = "alfred"
	app.Version = version
	app.Commands = cmd.Command()

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
