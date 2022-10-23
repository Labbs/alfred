package bookmark

import (
	"github.com/gofiber/fiber/v2"
	b "github.com/labbs/alfred/pkg/services/bookmark"
)

type bookmarkHander struct {
	bookmkark b.BookmarkRepository
}

func InitRoute(r fiber.Router) {
	h := bookmarkHander{bookmkark: b.NewBookmarkRepository()}

	g := r.Group("/bookmark")
	g.Post("/Create")
}
