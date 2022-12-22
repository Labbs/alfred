package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/logger"
	u "github.com/labbs/alfred/pkg/services/user"
)

type authHandler struct {
	user u.UserRepository
}

func InitRoute(r fiber.Router) {
	logger.Logger.Info().Msg("Initializing auth routes")
	h := authHandler{user: u.NewUserRepository()}

	r.Get("/auth/login", h.login)
	r.Post("/auth/login", h.login)
	r.Get("/auth/logout", h.logout)
}
