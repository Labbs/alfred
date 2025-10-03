package v1

import (
	"github.com/labbs/alfred/infrastructure"
	"github.com/labbs/alfred/interfaces/http/v1/auth"
)

func SetupRouterV1(deps infrastructure.Deps) {
	deps.Logger.Info().Str("component", "http.router.v1").Msg("Setting up API v1 routes")
	grp := deps.Http.FiberOapi.Group("/api/v1")

	authCtrl := auth.Controller{
		Config:    deps.Config,
		Logger:    deps.Logger,
		FiberOapi: grp.Group("/auth"),
	}
	auth.SetupAuthRouter(authCtrl)
}
