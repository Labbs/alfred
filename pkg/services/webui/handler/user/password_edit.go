package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func (h userHandler) passwordEdit(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
	c.ClearCookie("Error", "Success")

	user, err := h.user.FindUserByUsername(store.Get("username").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.update_avatar").Msg("could_not_get_user_profile")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get user profile"})
		return c.Redirect("/user/profile")
	}

	if c.FormValue("new-password") != c.FormValue("confirm-password") {
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Passwords do not match"})
		return c.Redirect("/user/profile")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.FormValue("current-password"))); err != nil {
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Current password is incorrect"})
		return c.Redirect("/user/profile")
	}

	hash, errB := bcrypt.GenerateFromPassword([]byte(c.FormValue("new-password")), 14)
	if errB != nil {
		logger.Logger.Error().Err(errB).Str("event", "user.update_avatar").Msg("could_not_update_user_avatar")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not update user avatar"})
		return c.Redirect("/user/profile")
	}

	user.Password = string(hash)

	err = h.user.UpdateUser(user)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.update_avatar").Msg("could_not_update_user_avatar")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not update user avatar"})
		return c.Redirect("/user/profile")
	}

	c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Password updated successfully"})
	return c.Redirect("/user/profile")
}
