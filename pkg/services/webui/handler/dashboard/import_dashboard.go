package dasbboard

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
	"github.com/labbs/alfred/pkg/services/dashboard"
)

func (h dashboardHandler) importDashboard(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
	c.ClearCookie("Error", "Success")

	form, errF := c.MultipartForm()
	if errF != nil {
		logger.Logger.Error().Err(errF).Str("event", "dashboard.import_dashboard").Msg("could_not_get_file")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get file"})
		return c.Redirect("/dashboard/list")
	}

	e := c.SaveFileToStorage(form.File["dashboard"][0], "dashboard", temporaryStore)
	if e != nil {
		logger.Logger.Error().Err(e).Str("event", "dashboard.import_dashboard").Msg("could_not_save_file")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not save file"})
		return c.Redirect("/dashboard/list")
	}

	file, errS := temporaryStore.Get("dashboard")
	if errS != nil {
		logger.Logger.Error().Err(errS).Str("event", "dashboard.import_dashboard").Msg("could_not_get_file_from_memory_storage")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get file from memory storage"})
		return c.Redirect("/dashboard/list")
	}

	var d dashboard.Dashboard

	err := json.Unmarshal(file, &d)
	if err != nil {
		logger.Logger.Error().Err(err).Str("event", "dashboard.import_dashboard").Msg("could_not_unmarshal_file")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not unmarshal file"})
		return c.Redirect("/dashboard/list")
	}

	d.UserId = store.Get("user_id").(string)
	d.Id = slug.Make(d.Name) + "-" + string(common.RandomString(5))

	errC := h.dashboard.CreateDashboard(d)
	if errC != nil {
		logger.Logger.Error().Err(errC.Error).Msg("Failed to create dashboard")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to create dashboard"})
		return c.Redirect("/dashboard/list")
	} else {
		c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Dashboard created successfully"})
		return c.Redirect("/dashboard/list")
	}
}
