package cmd

import "github.com/urfave/cli/v2"

func Command() []*cli.Command {
	return []*cli.Command{
		server(),
		migrate_database(),
		user(),
	}
}
