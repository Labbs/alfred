package user

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory"
	u "github.com/labbs/alfred/pkg/services/user"
)

type userHandler struct {
	user     u.UserRepository
	sessions *session.Store
}

var (
	temporaryStore *memory.Storage
)

func InitRoute(r fiber.Router, sessions *session.Store) {
	h := userHandler{user: u.NewUserRepository(), sessions: sessions}

	temporaryStore = memory.New(memory.Config{
		GCInterval: 10 * time.Second,
	})

	g := r.Group("/user")
	g.Get("/profile", h.profile)
	g.Post("/avatar/edit", h.avatarEdit)
	g.Get("/avatar/delete", h.avatarDelete)
	g.Post("/password/edit", h.passwordEdit)
	g.Post("/email/edit", h.emailEdit)
	g.Get("/light", h.lightDarkEdit)
	g.Get("/dark", h.lightDarkEdit)
	g.Post("/token/create", h.tokenCreate)
	g.Get("/token/delete/:id", h.tokenDelete)
}
