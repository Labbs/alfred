package webui

import (
	"github.com/gofiber/fiber/v2"
)

func checkSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		store, _ := sessions.Get(c)
		if c.Path() != "/auth/login" {
			if store.Get("username") == nil {
				return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
			}
		}
		return c.Next()
	}
}
