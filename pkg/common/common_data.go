package common

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/labbs/alfred/pkg/logger"
	"github.com/labbs/alfred/pkg/services/user"
)

func CommonData(sess *session.Store, c *fiber.Ctx) (fiber.Map, *session.Session) {
	d := make(fiber.Map)
	store, _ := sess.Get(c)
	r := user.NewUserRepository()
	u, err := r.FindUserByUsername(store.Get("username").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "common_data").Msg("failed to find user")
	}
	u.Password = ""
	d["Profile"] = u
	d["Avatar"] = template.URL(u.Avatar)
	d["LightDark"] = u.LightDark
	return d, store
}
