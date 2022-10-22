package dasbboard

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
	dash "github.com/labbs/alfred/pkg/services/dashboard"
)

func (h dashboardHandler) index(c *fiber.Ctx) error {
	d, store := common.CommonData(h.sessions, c)
	d["Title"] = "Dashboard"
	d["Error"] = c.Cookies("error-flash")
	d["Success"] = c.Cookies("success-flash")
	c.ClearCookie("error-flash", "success-flash")

	dashboard, err := h.dashboard.GetDefaultDashboard(store.Get("user_id").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "dashboard.get_dashboard").Msg("could_not_get_dashboard")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get default dashboard"})
		d["Dashboard"] = dash.Dashboard{}
	} else {
		d["Dashboard"] = dashboard
	}

	var css string
	var js string

	for _, widget := range dashboard.Widgets {
		css += widget.CSS
		js += widget.JS
	}

	d["CSS"] = template.CSS(css)
	d["JS"] = template.JS(js)

	return c.Render("templates/dashboard", d, "templates/layouts/main")
}
