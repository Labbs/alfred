package bookmark

import (
	fiber "github.com/gofiber/fiber/v2"
	b "github.com/labbs/alfred/pkg/services/bookmark"
	u "github.com/labbs/alfred/pkg/services/user"
)

type bookmarkHander struct {
	bookmark b.BookmarkRepository
	user     u.UserRepository
}

func InitRoute(r fiber.Router) {
	h := bookmarkHander{bookmark: b.NewBookmarkRepository(), user: u.NewUserRepository()}

	g := r.Group("/bookmark")
	g.Post("/create", h.createBookmark)
}
