package user

import (
	"encoding/base64"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
)

func (h userHandler) avatarEdit(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
	c.ClearCookie("Error", "Success")

	user, err := h.user.FindUserByUsername(store.Get("username").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.update_avatar").Msg("could_not_get_user_profile")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get user profile"})
		return c.Redirect("/user/profile")
	}

	form, errf := c.MultipartForm()
	if errf != nil {
		logger.Logger.Error().Err(errf).Str("event", "user.update_avatar").Msg("could_not_get_avatar_file")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get avatar file"})
		return c.Redirect("/user/profile")
	}

	e := c.SaveFileToStorage(form.File["avatar"][0], "avatars", temporaryStore)
	if e != nil {
		logger.Logger.Error().Err(e).Str("event", "user.update_avatar").Msg("could_not_save_avatar_file")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not save avatar file"})
		return c.Redirect("/user/profile")
	}

	file, errS := temporaryStore.Get("avatars")
	if errS != nil {
		logger.Logger.Error().Err(errS).Str("event", "user.update_avatar").Msg("could_not_get_avatar_file_from_memory_storage")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get avatar file from memory storage"})
		return c.Redirect("/user/profile")
	}

	var base64Encoding string

	switch form.File["avatar"][0].Header["Content-Type"][0] {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	base64Encoding += base64.StdEncoding.EncodeToString(file)

	user.Avatar = base64Encoding

	err = h.user.UpdateUser(user)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.update_avatar").Msg("could_not_update_user_avatar")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not update user avatar"})
		return c.Redirect("/user/profile")
	}

	c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Avatar updated successfully"})
	return c.Redirect("/user/profile")
}
