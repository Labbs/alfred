package dashboard

import (
	"encoding/json"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gosimple/slug"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
)

// index is the handler for the dashboard index page
// path: /
func (h dashboardHandler) index(c *fiber.Ctx) error {
	d := c.Locals("commonData").(fiber.Map)
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	d["Title"] = "Dashboard"
	d["Error"] = c.Cookies("error-flash")
	d["Success"] = c.Cookies("success-flash")
	c.ClearCookie("error-flash", "success-flash")

	logger.Logger.Debug().Str("event", "dashboard.get_dashboard").Msg("getting_dashboard")
	dashboard, err := h.dashboard.GetDefaultDashboard(store.Get("user_id").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "dashboard.get_dashboard").Msg("could_not_get_dashboard")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get default dashboard"})
		d["Dashboard"] = Dashboard{}
	} else {
		d["Dashboard"] = dashboard
	}

	logger.Logger.Debug().Str("event", "dashboard.get_dashboard").Msg("getting_dashboard_widgets")
	var css string
	var js string

	for _, widget := range dashboard.Widgets {
		css += widget.CSS
		js += widget.JS
	}

	d["CSS"] = template.CSS(css)
	d["JS"] = template.JS(js)

	logger.Logger.Debug().Str("event", "dashboard.get_dashboard").Msg("rendering_dashboard")
	return c.Render("templates/dashboard", d, "templates/layouts/main")
}

func (h dashboardHandler) createDashboard(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)

	var d Dashboard
	d.Name = c.FormValue("name")
	d.UserId = store.Get("user_id").(string)
	d.Id = slug.Make(d.Name) + "-" + string(common.RandomString(5))

	err := h.dashboard.CreateDashboard(d)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Msg("Failed to create dashboard")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to create dashboard"})
		return c.Redirect("/dashboard/list")
	} else {
		c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Dashboard created successfully"})
		return c.Redirect("/dashboard/list")
	}
}

func (h dashboardHandler) deleteDashboard(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)

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

func (h dashboardHandler) editDashboard(c *fiber.Ctx) error {
	d := c.Locals("commonData").(fiber.Map)
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
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

func (h dashboardHandler) editWidget(c *fiber.Ctx) error {
	d := c.Locals("commonData").(fiber.Map)
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
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

func (h dashboardHandler) exportDashboard(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	c.ClearCookie("error-flash", "success-flash")

	d, err := h.dashboard.GetDashboardById(store.Get("user_id").(string), c.Params("id"))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "dashboard.export_dashboard.get_dashboard").Msg("could_not_get_dashboard")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get dashboard"})
		return c.Redirect("/dashboard/list")
	}

	d.Id = ""
	d.Default = false

	var widgets []Widget
	for _, widget := range d.Widgets {
		widget.Id = ""
		widgets = append(widgets, widget)
	}

	d.Widgets = widgets

	c.Set("Content-Disposition", "attachment; filename="+d.Name+".json")
	return c.JSON(d)
}

func (h dashboardHandler) importDashboard(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
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

	var d Dashboard

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

func (h dashboardHandler) listDashboard(c *fiber.Ctx) error {
	d := c.Locals("commonData").(fiber.Map)
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
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

func (h dashboardHandler) saveDashboard(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	c.ClearCookie("Error", "Success")

	var widgets []Widget
	if err := c.BodyParser(&widgets); err != nil {
		logger.Logger.Error().Err(err).Str("event", "webui.updateDashboard").Msg("failed to parse dashboard")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to parse dashboard"})
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": true})
	}

	var inWidgets []Widget
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

func (h dashboardHandler) saveWidget(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)

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

func (h dashboardHandler) setDefaultDashboard(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)

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
