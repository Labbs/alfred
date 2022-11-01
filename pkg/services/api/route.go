package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/labbs/alfred/pkg/services/api/handler/bookmark"
	"github.com/labbs/alfred/pkg/services/user"
)

var (
	userRepository user.UserRepository
)

func InitRoute(r *fiber.App) {
	userRepository = user.NewUserRepository()
	g := r.Group("/api")

	// add cors middleware
	g.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	g.Use(checkToken())

	bookmark.InitRoute(g)
}
