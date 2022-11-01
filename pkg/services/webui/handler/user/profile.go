package user

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
)

func (h userHandler) profile(c *fiber.Ctx) error {
	d, _ := common.CommonData(h.sessions, c)
	d["Title"] = "Profile"
	d["Error"] = c.Cookies("error-flash")
	d["Success"] = c.Cookies("success-flash")
	c.ClearCookie("error-flash", "success-flash")
	return c.Render("templates/profile", d, "templates/layouts/main")
}
