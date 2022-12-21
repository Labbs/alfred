package user

import (
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func (h userHandler) profile(c *fiber.Ctx) error {
	d := c.Locals("commonData").(fiber.Map)
	d["Title"] = "Profile"
	d["Error"] = c.Cookies("error-flash")
	d["Success"] = c.Cookies("success-flash")
	c.ClearCookie("error-flash", "success-flash")
	return c.Render("templates/profile", d, "templates/layouts/main")
}

func (h userHandler) emailEdit(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	c.ClearCookie("Error", "Success")

	user, err := h.user.FindUserByUsername(store.Get("username").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.update_avatar").Msg("could_not_get_user_profile")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get user profile"})
		return c.Redirect("/user/profile")
	}

	user.Email = c.FormValue("email")

	err = h.user.UpdateUser(user)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.update_avatar").Msg("could_not_update_user_avatar")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not update user avatar"})
		return c.Redirect("/user/profile")
	}

	c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Email updated successfully"})
	return c.Redirect("/user/profile")
}

func (h userHandler) avatarDelete(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	c.ClearCookie("Error", "Success")

	user, err := h.user.FindUserByUsername(store.Get("username").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.update_avatar").Msg("could_not_get_user_profile")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get user profile"})
		return c.Redirect("/user/profile")
	}

	user.Avatar = ""

	err = h.user.UpdateUser(user)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.update_avatar").Msg("could_not_update_user_avatar")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not update user avatar"})
		return c.Redirect("/user/profile")
	}

	c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Avatar updated successfully"})
	return c.Redirect("/user/profile")
}

func (h userHandler) avatarEdit(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
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

func (h userHandler) darkMode(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	c.ClearCookie("Error", "Success")

	user, err := h.user.FindUserByUsername(store.Get("username").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.update_light_dark").Msg("could_not_get_user_profile")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get user profile"})
		return c.Redirect("/user/profile")
	}

	user.DarkMode = c.FormValue("dark_mode")

	err = h.user.UpdateUser(user)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.update_light_dark").Msg("could_not_update_user_light_dark")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not update user light/dark"})
		return c.Redirect("/user/profile")
	}

	return c.Redirect("/user/profile")
}

func (h userHandler) passwordEdit(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
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

func (h userHandler) tokenCreate(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	c.ClearCookie("error-flash", "success-flash")

	var token Token
	token.Name = c.FormValue("name")
	token.UserId = store.Get("user_id").(string)
	token.Id = string(common.RandomString(32))

	var scopes []Scope

	if c.FormValue("bookmark") == "on" {
		scopes = append(scopes, Scope{Name: "bookmark"})
	}

	token.Scope = scopes

	err := h.user.CreateToken(token)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.create_token").Msg("could_not_create_token")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not create token"})
	}

	return c.Redirect("/user/profile")
}

func (h userHandler) tokenDelete(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	c.ClearCookie("error-flash", "success-flash")

	err := h.user.DeleteTokenById(c.Params("id"), store.Get("user_id").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.delete_token").Msg("could_not_delete_token")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not delete token"})
	}

	logger.Logger.Info().Str("event", "user.delete_token").Msg("token_deleted")
	c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Token deleted"})
	return c.Redirect("/user/profile")
}
