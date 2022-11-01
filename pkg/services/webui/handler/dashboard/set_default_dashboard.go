package dasbboard

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
)

func (h dashboardHandler) setDefaultDashboard(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)

	err := h.dashboard.SetDefaultDashboard(c.Params("id"), store.Get("user_id").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Msg("Failed to set default dashboard")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to set default dashboard"})
		return c.Redirect("/dashboard/list")
	} else {
		c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Dashboard set as default successfully"})
		return c.Redirect("/dashboard/list")
	}
}
