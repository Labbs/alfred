package dasbboard

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"

	dash "github.com/labbs/alfred/pkg/services/dashboard"
)

func (h dashboardHandler) saveDashboard(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
	c.ClearCookie("Error", "Success")

	var widgets []dash.Widget
	if err := c.BodyParser(&widgets); err != nil {
		logger.Logger.Error().Err(err).Str("event", "webui.updateDashboard").Msg("failed to parse dashboard")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to parse dashboard"})
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": true})
	}

	var inWidgets []dash.Widget
	for _, widget := range widgets {
		if widget.Id == "" {
			widget.Id = utils.UUIDv4()
			widget.DashboardId = c.Params("id")
			widget.UserId = store.Get("user_id").(string)
			if err := h.dashboard.CreateWidget(widget); err != nil {
				logger.Logger.Error().Err(err.Error).Str("event", "dashboard.update.create_widget").Msg("could_not_create_widget")
				c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not create widget"})
				return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": true})
			}
			inWidgets = append(inWidgets, widget)
		} else {
			widget.DashboardId = c.Params("id")
			widget.UserId = store.Get("user_id").(string)
			if err := h.dashboard.UpdateWidget(widget); err != nil {
				logger.Logger.Error().Err(err.Error).Str("event", "dashboard.update.update_widget").Msg("could_not_update_widget")
				c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not update widget"})
				return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": true})
			}
			inWidgets = append(inWidgets, widget)
		}
	}

	savedWidgets, err := h.dashboard.GetWidgetsByDashboardId(c.Params("id"), store.Get("user_id").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "dashboard.update.get_widgets").Msg("could_not_get_widgets")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get widgets"})
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": true})
	}

	for _, widget := range savedWidgets {
		found := false
		for _, w := range inWidgets {
			if widget.Id == w.Id {
				found = true
			}
		}
		if !found {
			if err := h.dashboard.DeleteWidget(widget.Id, store.Get("user_id").(string)); err != nil {
				logger.Logger.Error().Err(err.Error).Str("event", "dashboard.update.delete_widget").Msg("could_not_delete_widget")
				c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not delete widget"})
				return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": true})
			}
		}
	}

	c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Dashboard updated successfully"})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
}
