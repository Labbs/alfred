package auth

import "github.com/gofiber/fiber/v2"

func (h authHandler) logout(c *fiber.Ctx) error {
	store, _ := h.sessions.Get(c)
	store.Delete("user_id")
	store.Delete("username")
	store.Destroy()
	return c.Redirect("/auth/login")
}
