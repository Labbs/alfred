package application

import (
	"github.com/labbs/alfred/domain"
	"github.com/labbs/alfred/infrastructure/config"
	"github.com/rs/zerolog"
)

type UserApp struct {
	Config   config.Config
	Logger   zerolog.Logger
	UserPres domain.UserPers
}

func NewUserApp(config config.Config, logger zerolog.Logger, userPers domain.UserPers) *UserApp {
	return &UserApp{
		Config:   config,
		Logger:   logger,
		UserPres: userPers,
	}
}
