package dashboard

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory"
	"github.com/labbs/alfred/pkg/logger"
)

type dashboardHandler struct {
	dashboard DashboardRepository
}

var (
	temporaryStore *memory.Storage
)

func InitRoute(r fiber.Router) {
	logger.Logger.Info().Msg("Initializing dashboard routes")

	// initialize dashboard handler
	h := dashboardHandler{dashboard: NewDashboardRepository()}

	// initialize temporary storage
	temporaryStore = memory.New(memory.Config{
		GCInterval: 10 * time.Second,
	})

	// initialize dashboard routes
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
	g.Get("/view/:id", h.index)
}
