package auth

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/logger"
)

func (h authHandler) logout(c *fiber.Ctx) error {
	store, _ := h.sessions.Get(c)
	store.Delete("user_id")
	store.Delete("username")
	err := store.Destroy()
	if err != nil {
		logger.Logger.Error().Err(err).Str("event", "auth.logout").Msg(err.Error())
	}
	return c.Redirect("/auth/login")
}
