package common

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/labbs/alfred/pkg/services/user"
)

func CommonData(sess *session.Store, c *fiber.Ctx) (fiber.Map, *session.Session) {
	d := make(fiber.Map)
	store, _ := sess.Get(c)
	r := user.NewUserRepository()
	u, _ := r.FindUserByUsername(store.Get("username").(string))
	u.Password = ""
	d["Profile"] = u
	d["Avatar"] = template.URL(u.Avatar)
	return d, store
}
