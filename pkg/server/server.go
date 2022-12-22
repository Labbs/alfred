package server

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/urfave/cli/v2"

	"github.com/labbs/alfred/pkg/common/gorm"
	"github.com/labbs/alfred/pkg/config"
	"github.com/labbs/alfred/pkg/database"
	"github.com/labbs/alfred/pkg/logger"
	"github.com/labbs/alfred/pkg/services/auth"
	"github.com/labbs/alfred/pkg/services/bookmark"
	"github.com/labbs/alfred/pkg/services/dashboard"
	"github.com/labbs/alfred/pkg/services/user"
)

func RunServer(ctx *cli.Context) error {
	// Init logger
	logger.InitLogger(ctx.App.Version)

	// Init database
	database.InitDatabase()

	config.Version = ctx.App.Version

	// Init cron scheduler
	// cronscheduler.InitScheduler()

	r := fiber.New(fiber.Config{
		Views:                 engineInit(),
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			c.Cookie(&fiber.Cookie{Name: "error-flash", Value: err.Error()})
			return c.Redirect("/")
		},
	})

	// init sessions storage
	storage := gorm.New(gorm.Config{
		DB:    database.GetDbConnection().DB,
		Table: "session",
	})
	sessions := session.New(session.Config{
		Storage: storage,
	})
	// register types to be stored in the session
	sessions.RegisterType(map[string]interface{}{})
	sessions.RegisterType([]interface{}{})

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

	// enable static files
	r.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirStatic),
		PathPrefix: "static",
		Browse:     true,
	}))

	// enable the custom middleware to set the sessions in the context
	r.Use(func(c *fiber.Ctx) error {
		c.Locals("sessions", sessions)
		return c.Next()
	})

	// enable the custom middleware to check if the user is logged in
	r.Use(checkSession())

	r.Use(func(c *fiber.Ctx) error {
		if c.Path() != "/auth/login" {
			c.Locals("commonData", commonDataResponse(c))
		}
		return c.Next()
	})

	// enable the routes
	auth.InitRoute(r)
	user.InitRoute(r)
	bookmark.InitRoute(r)
	dashboard.InitRoute(r)

	logger.Logger.Info().Msg("app listening on 0.0.0.0:" + strconv.Itoa(config.Port))
	return r.Listen(":" + strconv.Itoa(config.Port))
}
