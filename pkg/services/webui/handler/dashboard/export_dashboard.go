package dasbboard

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
	"github.com/labbs/alfred/pkg/services/dashboard"
)

func (h dashboardHandler) exportDashboard(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
	c.ClearCookie("error-flash", "success-flash")

	d, err := h.dashboard.GetDashboardById(store.Get("user_id").(string), c.Params("id"))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "dashboard.export_dashboard.get_dashboard").Msg("could_not_get_dashboard")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get dashboard"})
		return c.Redirect("/dashboard/list")
	}

	d.Id = ""
	d.Default = false

	var widgets []dashboard.Widget
	for _, widget := range d.Widgets {
		widget.Id = ""
		widgets = append(widgets, widget)
	}

	d.Widgets = widgets

	c.Set("Content-Disposition", "attachment; filename="+d.Name+".json")
	return c.JSON(d)
}
