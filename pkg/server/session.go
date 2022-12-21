package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/labbs/alfred/pkg/logger"
)

func checkSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		store, err := c.Locals("sessions").(*session.Store).Get(c)
		if err != nil {
			logger.Logger.Error().Err(err).Str("event", "webui.check_session").Msg("failed to get session")
		}
		if c.Path() != "/auth/login" {
			if store.Get("username") == nil {
				return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
			}
		}
		return c.Next()
	}
}
