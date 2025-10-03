package application

import (
	"github.com/labbs/alfred/domain"
	"github.com/labbs/alfred/infrastructure/config"
	"github.com/rs/zerolog"
)

type AuthApp struct {
	Config      config.Config
	Logger      zerolog.Logger
	UserPers    domain.UserPers
	SessionPers domain.SessionPers
}

func NewAuthApp(config config.Config, logger zerolog.Logger, userPers domain.UserPers, sessionPers domain.SessionPers) *AuthApp {
	return &AuthApp{
		Config:      config,
		Logger:      logger,
		UserPers:    userPers,
		SessionPers: sessionPers,
	}
}
