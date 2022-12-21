package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory"
	"github.com/labbs/alfred/pkg/logger"
)

type userHandler struct {
	user UserRepository
}

var (
	temporaryStore *memory.Storage
)

func InitRoute(r fiber.Router) {
	logger.Logger.Info().Msg("Initializing user routes")
	h := userHandler{user: NewUserRepository()}

	temporaryStore = memory.New(memory.Config{
		GCInterval: 10 * time.Second,
	})

	g := r.Group("/user")
	g.Get("/profile", h.profile)
	g.Post("/avatar/edit", h.avatarEdit)
	g.Get("/avatar/delete", h.avatarDelete)
	g.Post("/password/edit", h.passwordEdit)
	g.Post("/email/edit", h.emailEdit)
	g.Post("/dark_mode", h.darkMode)
	g.Post("/token/create", h.tokenCreate)
	g.Get("/token/delete/:id", h.tokenDelete)
}
