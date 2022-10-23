package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/services/user"
)

var (
	userRepository user.UserRepository
)

func InitRoute(r *fiber.App) {
	userRepository = user.NewUserRepository()
	g := r.Group("/api")
	g.Use(checkToken())
}
