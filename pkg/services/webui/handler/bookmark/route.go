package bookmark

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	b "github.com/labbs/alfred/pkg/services/bookmark"
)

type bookmarkHandler struct {
	// user     u.UserRepository
	bookmark b.BookmarkRepository
	sessions *session.Store
}

func InitRoute(r fiber.Router, sessions *session.Store) {
	h := bookmarkHandler{bookmark: b.NewBookmarkRepository(), sessions: sessions}
	g := r.Group("/bookmark")
	g.Get("/", h.bookmarkList)
	g.Post("/", h.bookmarkList)
	g.Post("/create", h.createBookmark)
	g.Post("/create/bulk", h.createBulkBookmark)
	g.Post("/edit/:id", h.editBookmark)
	g.Get("/delete/:id", h.deleteBookmark)
}
