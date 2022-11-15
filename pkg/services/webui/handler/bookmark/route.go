package bookmark

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory"
	b "github.com/labbs/alfred/pkg/services/bookmark"
)

type bookmarkHandler struct {
	// user     u.UserRepository
	bookmark b.BookmarkRepository
	sessions *session.Store
}

var (
	temporaryStore *memory.Storage
)

func InitRoute(r fiber.Router, sessions *session.Store) {
	h := bookmarkHandler{bookmark: b.NewBookmarkRepository(), sessions: sessions}

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
