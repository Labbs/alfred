package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func (h authHandler) login(c *fiber.Ctx) error {
	d := make(fiber.Map)
	d["Error"] = c.Cookies("error-flash")
	d["Success"] = c.Cookies("success-flash")
	c.ClearCookie("error-flash", "success-flash")

	if c.Method() == "POST" {

		user, errFind := h.user.FindUserByUsername(c.FormValue("username"))
		if errFind != nil {
			logger.Logger.Error().Err(errFind.Error).Str("event", "webui.login.find_user").Msg("failed to find user")
			d["Error"] = "Failed to login"
			return c.Render("templates/login", d)
		}

		if user.Username == "" {
			logger.Logger.Error().Err(fmt.Errorf("user not found")).Str("event", "webui.login.user_exist").Msg("user not exist: " + c.FormValue("username"))
			d["Error"] = "User not exist"
			return c.Render("templates/login", d)
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.FormValue("password")))
		if err != nil {
			logger.Logger.Error().Err(err).Str("event", "webui.login.compare_password").Msg("incorrect password for " + user.Username)
			d["Error"] = "Incorrect password"
			return c.Render("templates/login", d)
		}

		store, _ := h.sessions.Get(c)
		store.Set("username", user.Username)
		store.Set("user_id", user.Id)

		errStore := store.Save()
		if errStore != nil {
			logger.Logger.Error().Err(errStore).Str("event", "webui.login").Msg("failed to save session")
			d["Error"] = "Internal server error"
			return c.Render("templates/login", d)
		}

		return c.Redirect("/")
	}

	return c.Render("templates/login", d)
}
