package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/labbs/alfred/pkg/logger"
	u "github.com/labbs/alfred/pkg/services/user"
)

type authHandler struct {
	user     u.UserRepository
	sessions *session.Store
}

func InitRoute(r fiber.Router, sessions *session.Store) {
	logger.Logger.Info().Msg("Initializing auth routes")
	h := authHandler{user: u.NewUserRepository(), sessions: sessions}
	r.Get("/auth/login", h.login)
	r.Post("/auth/login", h.login)
	r.Get("/auth/logout", h.logout)
}
