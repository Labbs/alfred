package dasbboard

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
)

func (h dashboardHandler) deleteDashboard(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)

	err := h.dashboard.DeleteDashboard(c.Params("id"), store.Get("user_id").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Msg("Failed to delete dashboard")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to delete dashboard"})
		return c.Redirect("/dashboard/list")
	} else {
		c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Dashboard deleted successfully"})
		return c.Redirect("/dashboard/list")
	}
}
