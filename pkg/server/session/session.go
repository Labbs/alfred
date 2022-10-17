package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func CheckSession(sessions *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		store, _ := sessions.Get(c)
		if store.Get("user_id") == nil {
			if c.Path() == "/app" {
				return c.Redirect("/app/login", fiber.StatusTemporaryRedirect)
			}
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"message": "unauthorized",
			})
		}
		c.Next()
		return nil
	}
}
