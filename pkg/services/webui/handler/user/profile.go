package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
)

func (h userHandler) profile(c *fiber.Ctx) error {
	d, _ := common.CommonData(h.sessions, c)
	d["Title"] = "Profile"
	return c.Render("templates/profile", d, "templates/layouts/main")
}
