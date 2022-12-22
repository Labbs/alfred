package server

import (
	"html/template"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/labbs/alfred/pkg/logger"
	"github.com/labbs/alfred/pkg/services/user"
)

func commonDataResponse(c *fiber.Ctx) fiber.Map {
	d := make(fiber.Map)
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	r := user.NewUserRepository()
	u, err := r.FindUserByUsername(store.Get("username").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "common_data").Msg("failed to find user")
	}
	u.Password = ""
	d["Profile"] = u
	d["Avatar"] = template.URL(u.Avatar)
	d["DarkMode"] = u.DarkMode
	d["FullScreen"] = false
	if c.Query("fullscreen") == "true" {
		d["FullScreen"] = true
	}
	if c.Path() == "/" || strings.Contains("/dashboard/view/", c.Path()) {
		d["Dashboard"] = true
	}
	return d
}
