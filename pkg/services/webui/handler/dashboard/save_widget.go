package dasbboard

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
)

func (h dashboardHandler) saveWidget(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)

	widget, err := h.dashboard.GetWidgetById(c.Params("id"), store.Get("user_id").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "dashboard.save_widget.get_widget").Msg("could_not_get_widget")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get widget"})
		return c.Redirect("/dashboard/list")
	}

	widget.Name = c.FormValue("name")
	widget.HTML = c.FormValue("html")
	widget.JS = c.FormValue("js")
	widget.CSS = c.FormValue("css")

	err = h.dashboard.UpdateWidget(widget)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "dashboard.save_widget.update_widget").Msg("could_not_update_widget")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not update widget"})
		return c.Redirect("/dashboard/list")
	} else {
		c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Widget updated successfully"})
		return c.Redirect("/dashboard/edit/" + widget.DashboardId)
	}
}
