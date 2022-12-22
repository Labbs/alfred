package bookmark

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory"
	"github.com/labbs/alfred/pkg/logger"
)

type bookmarkHandler struct {
	// user     u.UserRepository
	bookmark BookmarkRepository
}

var (
	temporaryStore *memory.Storage
)

func InitRoute(r fiber.Router) {
	logger.Logger.Info().Msg("Initializing bookmark routes")
	h := bookmarkHandler{bookmark: NewBookmarkRepository()}

	temporaryStore = memory.New(memory.Config{
		GCInterval: 10 * time.Second,
	})

	g := r.Group("/bookmark")
	g.Get("/", h.bookmarkList)
	g.Post("/", h.bookmarkList)
	g.Post("/create", h.createBookmark)
	g.Post("/create/bulk", h.createBulkBookmark)
	g.Post("/edit/:id", h.editBookmark)
	g.Get("/delete/:id", h.deleteBookmark)
	g.Get("/tags/clean_unused", h.cleanUnusedTags)
	g.Post("/import", h.importBookmarksFile)
}
