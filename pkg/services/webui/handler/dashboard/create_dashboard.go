package dasbboard

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
	"github.com/labbs/alfred/pkg/services/dashboard"
)

func (h dashboardHandler) createDashboard(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)

	var d dashboard.Dashboard
	d.Name = c.FormValue("name")
	d.UserId = store.Get("user_id").(string)
	d.Id = slug.Make(d.Name) + "-" + string(common.RandomString(5))

	err := h.dashboard.CreateDashboard(d)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Msg("Failed to create dashboard")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to create dashboard"})
		return c.Redirect("/dashboard/list")
	} else {
		c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Dashboard created successfully"})
		return c.Redirect("/dashboard/list")
	}
}
