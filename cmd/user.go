package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2/utils"
	d "github.com/labbs/alfred/pkg/database"
	"github.com/labbs/alfred/pkg/logger"
	model "github.com/labbs/alfred/pkg/services/user"

	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	username string
)

func user() *cli.Command {
	return &cli.Command{
		Name:  "user",
		Flags: databaseFlags,
		Subcommands: []*cli.Command{
			{
				Name:   "list",
				Action: list,
			},
			{
				Name:   "delete",
				Action: delete,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "username",
						Destination: &username,
					},
				},
			},
			{
				Name:   "add",
				Action: add,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "username",
						Destination: &username,
					},
				},
			},
			{
				Name:   "reset",
				Action: reset,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "username",
						Destination: &username,
					},
				},
			},
		},
	}
}

func list(ctx *cli.Context) error {
	// Init logger
	logger.InitLogger(ctx.App.Version)

	// Init database
	d.InitDatabase()

	c := d.GetDbConnection()

	u := make([]model.User, 0)
	r := c.DB.Find(&u)
	if r.Error != nil {
		logger.Logger.Error().Err(r.Error).Str("event", "cmd.user.list").Msg(r.Error.Error())
		return r.Error
	}
	for _, user := range u {
		fmt.Println(user.Username)
	}
	return nil
}

func delete(ctx *cli.Context) error {
	// Init logger
	logger.InitLogger(ctx.App.Version)

	if username == "" {
		logger.Logger.Error().Str("event", "cmd.user.delete").Msg("missing username")
		return nil
	}

	// Init database
	d.InitDatabase()

	c := d.GetDbConnection()

	u := model.User{Username: username}
	r := c.DB.Where("username = ?", username).Find(&u)
	if r.Error != nil {
		logger.Logger.Error().Err(r.Error).Str("event", "cmd.user.add").Str("username", username).Msg(r.Error.Error())
		return r.Error
	}

	if u.Id == "" {
		logger.Logger.Error().Str("event", "cmd.user.delete").Str("username", username).Msg("user not exist")
		return nil
	}

	rr := c.DB.Where("username = ?", username).Delete(&u)
	if r.Error != nil {
		logger.Logger.Error().Err(rr.Error).Str("event", "cmd.user.delete").Str("username", username).Msg(rr.Error.Error())
		return rr.Error
	}
	return nil
}

func add(ctx *cli.Context) error {
	// Init logger
	logger.InitLogger(ctx.App.Version)

	if username == "" {
		logger.Logger.Error().Str("event", "cmd.user.add").Msg("missing username")
		return nil
	}

	// Init database
	d.InitDatabase()

	c := d.GetDbConnection()

	u := model.User{Username: username}

	r := c.DB.Where("username = ?", u.Username).Find(&u)
	if r.Error != nil {
		logger.Logger.Error().Err(r.Error).Str("event", "cmd.user.add").Str("username", username).Msg(r.Error.Error())
		return r.Error
	}

	if u.Id != "" {
		logger.Logger.Error().Str("event", "cmd.user.add").Str("username", username).Msg("user already exist")
		return nil
	}

	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 30)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	bytes, err := bcrypt.GenerateFromPassword(b, 14)
	if err != nil {
		logger.Logger.Error().Err(err).Str("event", "cmd.user.add.generate_from_password").Str("username", username).Msg(err.Error())
		return err
	}

	u.Id = utils.UUIDv4()
	u.Password = string(bytes)

	rr := c.DB.Create(&u)
	if rr.Error != nil {
		logger.Logger.Error().Err(rr.Error).Str("event", "cmd.user.add").Str("username", username).Msg(rr.Error.Error())
		return rr.Error
	}

	fmt.Println("username: " + username)
	fmt.Println("password: " + string(b))

	return nil
}

func reset(ctx *cli.Context) error {
	// Init logger
	logger.InitLogger(ctx.App.Version)

	if username == "" {
		logger.Logger.Error().Str("event", "cmd.user.reset").Msg("missing username")
		return nil
	}

	// Init database
	d.InitDatabase()

	c := d.GetDbConnection()

	u := model.User{Username: username}
	r := c.DB.Where("username = ?", u.Username).Find(&u)
	if r.Error != nil {
		logger.Logger.Error().Err(r.Error).Str("event", "cmd.user.reset").Str("username", username).Msg(r.Error.Error())
		return r.Error
	}

	if u.Id == "" {
		logger.Logger.Error().Str("event", "cmd.user.reset").Str("username", username).Msg("user not exist")
		return nil
	}

	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 30)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	bytes, err := bcrypt.GenerateFromPassword(b, 14)
	if err != nil {
		logger.Logger.Error().Err(err).Str("event", "cmd.user.reset.generate_from_password").Str("username", username).Msg(err.Error())
		return err
	}

	u.Password = string(bytes)

	rr := c.DB.Save(&u)
	if rr.Error != nil {
		logger.Logger.Error().Err(rr.Error).Str("event", "cmd.user.reset").Str("username", username).Msg(rr.Error.Error())
		return rr.Error
	}

	fmt.Println("username: " + username)
	fmt.Println("password: " + string(b))

	return nil
}
