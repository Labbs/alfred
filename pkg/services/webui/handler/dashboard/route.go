package dasbboard

import (
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/labbs/alfred/pkg/config"
	"github.com/labbs/alfred/pkg/logger"
	d "github.com/labbs/alfred/pkg/services/dashboard"
	"gopkg.in/yaml.v2"
)

type dashboardHandler struct {
	dashboard   d.DashboardRepository
	sessions    *session.Store
	widgetsList []d.WidgetType
}

func InitRoute(r fiber.Router, sessions *session.Store) {
	h := dashboardHandler{dashboard: d.NewDashboardRepository(), sessions: sessions}

	if _, err := os.Stat(config.WidgetsList); err == nil {
		file, _ := ioutil.ReadFile(config.WidgetsList)
		err = yaml.Unmarshal(file, &h.widgetsList)
		if err != nil {
			logger.Logger.Error().Err(err).Str("event", "dashboard.init_route").Msg("could_not_unmarshal_widgets_list")
			h.widgetsList = []d.WidgetType{}
		}
	} else {
		logger.Logger.Error().Err(err).Str("event", "dashboard.init_route").Msg("widgets_list_file_not_found")
		h.widgetsList = []d.WidgetType{}
	}

	r.Get("/", h.index)

	g := r.Group("/dashboard")
	g.Get("/list", h.listDashboard)
	g.Post("/create", h.createDashboard)
	g.Get("/delete/:id", h.deleteDashboard)
	g.Get("/set_default/:id", h.setDefaultDashboard)
	g.Get("/edit/:id", h.editDashboard)
	g.Post("/save/:id", h.saveDashboard)
}
