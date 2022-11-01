package dasbboard

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
)

func (h dashboardHandler) editDashboard(c *fiber.Ctx) error {
	d, store := common.CommonData(h.sessions, c)
	d["Title"] = "Dashboard"
	d["Error"] = c.Cookies("error-flash")
	d["Success"] = c.Cookies("success-flash")
	c.ClearCookie("error-flash", "success-flash")

	dashboard, err := h.dashboard.GetDashboardById(store.Get("user_id").(string), c.Params("id"))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "dashboard.get_dashboard").Msg("could_not_get_dashboard")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get dashboard"})
		return c.Redirect("/dashboard/list")
	}

	d["Dashboard"] = dashboard

	return c.Render("templates/dashboard_edit", d, "templates/layouts/main")
}
