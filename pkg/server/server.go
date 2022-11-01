package server

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/urfave/cli/v2"

	"github.com/labbs/alfred/pkg/config"
	"github.com/labbs/alfred/pkg/database"
	"github.com/labbs/alfred/pkg/logger"
	"github.com/labbs/alfred/pkg/services/api"
	"github.com/labbs/alfred/pkg/services/cronscheduler"
	"github.com/labbs/alfred/pkg/services/webui"
)

func RunServer(ctx *cli.Context) error {
	// Init logger
	logger.InitLogger(ctx.App.Version)

	// Init database
	database.InitDatabase()

	config.Version = ctx.App.Version

	// Init cron scheduler
	cronscheduler.InitScheduler()

	r := fiber.New(fiber.Config{
		Views:                 webui.EngineInit(),
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			c.Cookie(&fiber.Cookie{Name: "error-flash", Value: err.Error()})
			return c.Redirect("/")
		},
	})

	// enable gofiber logs (custom middleware)
	if config.EnableHTTPLogs {
		r.Use(logger.New())
	}

	// gofiber recover => https://docs.gofiber.io/api/middleware/recover
	r.Use(recover.New())

	// enable monitor page => https://docs.gofiber.io/api/middleware/monitor
	if config.EnableMetricsPage {
		r.Get("/metrics", monitor.New(monitor.Config{Title: "Alfred Metrics Page"}))
	}

	api.InitRoute(r)
	// init WebUI routes
	webui.InitRoute(r)

	logger.Logger.Info().Msg("app listening on 0.0.0.0:" + strconv.Itoa(config.Port))
	return r.Listen(":" + strconv.Itoa(config.Port))
}
