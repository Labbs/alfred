package dasbboard

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory"
	d "github.com/labbs/alfred/pkg/services/dashboard"
)

type dashboardHandler struct {
	dashboard d.DashboardRepository
	sessions  *session.Store
}

var (
	temporaryStore *memory.Storage
)

func InitRoute(r fiber.Router, sessions *session.Store) {
	h := dashboardHandler{dashboard: d.NewDashboardRepository(), sessions: sessions}

	temporaryStore = memory.New(memory.Config{
		GCInterval: 10 * time.Second,
	})

	r.Get("/", h.index)

	g := r.Group("/dashboard")
	g.Get("/list", h.listDashboard)
	g.Post("/create", h.createDashboard)
	g.Get("/delete/:id", h.deleteDashboard)
	g.Get("/set_default/:id", h.setDefaultDashboard)
	g.Get("/widget/edit/:id", h.editWidget)
	g.Post("/widget/save/:id", h.saveWidget)
	g.Get("/edit/:id", h.editDashboard)
	g.Post("/save/:id", h.saveDashboard)
	g.Get("/export/:id", h.exportDashboard)
	g.Post("/import", h.importDashboard)
}
