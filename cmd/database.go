package cmd

import (
	"github.com/labbs/alfred/pkg/database"
	"github.com/labbs/alfred/pkg/logger"
	b "github.com/labbs/alfred/pkg/services/bookmark"
	d "github.com/labbs/alfred/pkg/services/dashboard"
	u "github.com/labbs/alfred/pkg/services/user"
	"github.com/urfave/cli/v2"
)

func migrate_database() *cli.Command {
	return &cli.Command{
		Name:  "database",
		Flags: databaseFlags,
		Subcommands: []*cli.Command{
			{
				Name:   "migrate",
				Action: migrate,
			},
		},
	}
}

func migrate(ctx *cli.Context) error {
	// Init logger
	logger.InitLogger(ctx.App.Version)

	database.InitDatabase()

	db := database.GetDbConnection()

	err := db.DB.AutoMigrate(&u.User{}, &b.Bookmark{}, &b.Tag{}, &d.Dashboard{}, &d.Widget{}, &u.Token{})
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Error migrating database")
		return nil
	}
	return nil
}
