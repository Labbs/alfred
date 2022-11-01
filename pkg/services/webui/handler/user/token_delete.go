package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
)

func (h userHandler) tokenDelete(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
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
