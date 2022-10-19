package cmd

import (
	"github.com/labbs/alfred/pkg/config"
	"github.com/urfave/cli/v2"
)

var databaseFlags = []cli.Flag{
	&cli.StringFlag{
		Name:        "engine",
		Aliases:     []string{"e"},
		EnvVars:     []string{"ENGINE"},
		Value:       "sqlite3",
		Usage:       "Database engine",
		Destination: &config.Database.Engine,
	},
	&cli.StringFlag{
		Name:        "dsn",
		Aliases:     []string{"d"},
		EnvVars:     []string{"DSN"},
		Value:       "alfred.sqlite",
		Usage:       "Database DSN",
		Destination: &config.Database.DSN,
	},
	&cli.BoolFlag{
		Name:    "debug",
		Aliases: []string{"dg"},
		EnvVars: []string{"DEBUG"},
		Value:   false,
		Usage:   "Debug mode",
	},
}
