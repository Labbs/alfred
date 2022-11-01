package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
	"github.com/labbs/alfred/pkg/services/user"
)

func (h userHandler) tokenCreate(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
	c.ClearCookie("error-flash", "success-flash")

	var token user.Token
	token.Name = c.FormValue("name")
	token.UserId = store.Get("user_id").(string)
	token.Id = string(common.RandomString(32))

	var scopes []user.Scope

	if c.FormValue("bookmark") == "on" {
		scopes = append(scopes, user.Scope{Name: "bookmark"})
	}

	token.Scope = scopes

	err := h.user.CreateToken(token)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "user.create_token").Msg("could_not_create_token")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not create token"})
	}

	return c.Redirect("/user/profile")
}
