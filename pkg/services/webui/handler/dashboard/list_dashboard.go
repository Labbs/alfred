package dasbboard

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
)

func (h dashboardHandler) listDashboard(c *fiber.Ctx) error {
	d, store := common.CommonData(h.sessions, c)
	d["Title"] = "Dashboard list"
	d["Error"] = c.Cookies("error-flash")
	d["Success"] = c.Cookies("success-flash")
	c.ClearCookie("error-flash", "success-flash")

	dashboards, err := h.dashboard.GetAllDashboards(store.Get("user_id").(string))
	if err != nil {
		d["Error"] = "Failed to get dashboards"
		d["Dashboards"] = []string{}
	} else {
		if dashboards == nil {
			d["Dashboards"] = []string{}
		} else {
			d["Dashboards"] = dashboards
		}
	}

	return c.Render("templates/dashboard_list", d, "templates/layouts/main")
}
