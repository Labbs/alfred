package user

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
)

func (h userHandler) lightDarkEdit(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
	c.ClearCookie("Error", "Success")

	user, err := h.user.FindUserByUsername(store.Get("username").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.update_light_dark").Msg("could_not_get_user_profile")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get user profile"})
		return c.Redirect("/")
	}
	if c.Path() == "/user/light" {
		user.LightDark = "light"
	}
	if c.Path() == "/user/dark" {
		user.LightDark = "dark"
	}

	err = h.user.UpdateUser(user)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.update_light_dark").Msg("could_not_update_user_light_dark")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not update user light/dark"})
		return c.Redirect("/")
	}

	return c.Redirect("/")
}
