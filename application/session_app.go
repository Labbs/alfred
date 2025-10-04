package application

import (
	"github.com/labbs/alfred/domain"
	"github.com/labbs/alfred/infrastructure/config"
	"github.com/rs/zerolog"
)

type SessionApp struct {
	Config      config.Config
	Logger      zerolog.Logger
	SessionPers domain.SessionPers
}

func NewSessionApp(config config.Config, logger zerolog.Logger, sessionPers domain.SessionPers) *SessionApp {
	return &SessionApp{
		Config:      config,
		Logger:      logger,
		SessionPers: sessionPers,
	}
}
