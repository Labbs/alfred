package infrastructure

import (
	"github.com/labbs/alfred/application"
	"github.com/labbs/alfred/infrastructure/config"
	"github.com/labbs/alfred/infrastructure/cronscheduler"
	"github.com/labbs/alfred/infrastructure/database"
	"github.com/labbs/alfred/infrastructure/http"
	"github.com/rs/zerolog"
)

type Deps struct {
	Config        config.Config
	Logger        zerolog.Logger
	Http          http.Config
	CronScheduler cronscheduler.Config
	Database      database.Config

	UserApp    *application.UserApp
	SessionApp *application.SessionApp
	AuthApp    *application.AuthApp
}
