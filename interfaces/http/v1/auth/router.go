package auth

import (
	"github.com/labbs/alfred/infrastructure/config"
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/rs/zerolog"
)

type Controller struct {
	Config    config.Config
	Logger    zerolog.Logger
	FiberOapi *fiberoapi.OApiGroup
}

func SetupAuthRouter(ctrl Controller) {}
