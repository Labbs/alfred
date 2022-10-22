package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/labbs/alfred/pkg/config"
	s "github.com/labbs/alfred/pkg/server"
)

func server() *cli.Command {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:        "port",
			Aliases:     []string{"p"},
			EnvVars:     []string{"PORT"},
			Value:       8080,
			Usage:       "Server Port",
			Destination: &config.Port,
		},
		&cli.StringFlag{
			Name:        "session-secret-key",
			Aliases:     []string{"ssk"},
			EnvVars:     []string{"SESSION_SECRET_KEY"},
			Value:       "secret",
			Usage:       "Session secret key",
			Destination: &config.Session.SecretKey,
		},
		&cli.IntFlag{
			Name:        "session-expire",
			Aliases:     []string{"se"},
			EnvVars:     []string{"SESSION_EXPIRE"},
			Value:       3600,
			Usage:       "Session expire time",
			Destination: &config.Session.Expire,
		},
		&cli.BoolFlag{
			Name:        "enable-http-logs",
			Aliases:     []string{"ehl"},
			EnvVars:     []string{"ENABLE_HTTP_LOGS"},
			Value:       false,
			Usage:       "Enable http server logs",
			Destination: &config.EnableHTTPLogs,
		},
		&cli.BoolFlag{
			Name:        "enable-metrics-page",
			Aliases:     []string{"emp"},
			EnvVars:     []string{"ENABLE_METRICS_PAGE"},
			Value:       false,
			Usage:       "Enable metrics page",
			Destination: &config.EnableMetricsPage,
		},
	}
	flags = append(flags, databaseFlags...)
	return &cli.Command{
		Name:   "server",
		Action: s.RunServer,
		Flags:  flags,
	}
}
