package dasbboard

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	d "github.com/labbs/alfred/pkg/services/dashboard"
)

type dashboardHandler struct {
	dashboard d.DashboardRepository
	sessions  *session.Store
}

func InitRoute(r fiber.Router, sessions *session.Store) {
	h := dashboardHandler{dashboard: d.NewDashboardRepository(), sessions: sessions}

	r.Get("/", h.index)

	g := r.Group("/dashboard")
	g.Get("/list", h.listDashboard)
	g.Post("/create", h.createDashboard)
	g.Get("/delete/:id", h.deleteDashboard)
	g.Get("/set_default/:id", h.setDefaultDashboard)
	g.Get("/edit/:id", h.editDashboard)
	g.Post("/save/:id", h.saveDashboard)
}
