package dasbboard

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
)

func (h dashboardHandler) editWidget(c *fiber.Ctx) error {
	d, store := common.CommonData(h.sessions, c)
	d["Title"] = "Dashboard"
	d["Error"] = c.Cookies("error-flash")
	d["Success"] = c.Cookies("success-flash")
	c.ClearCookie("error-flash", "success-flash")

	widget, err := h.dashboard.GetWidgetById(c.Params("id"), store.Get("user_id").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "dashboard.edit_widget.get_widget").Msg("could_not_get_widget")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get widget"})
		return c.Redirect("/dashboard/list")
	}

	d["Widget"] = widget
	return c.Render("templates/widget_edit", d, "templates/layouts/main")
}
