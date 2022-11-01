package bookmark

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
)

func (h bookmarkHandler) deleteBookmark(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
	c.ClearCookie("Error", "Success")

	err := h.bookmark.DeleteBookmark(c.Params("id"), store.Get("user_id").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "webui.deleteBookmark").Msg("failed to delete bookmark")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to delete bookmark"})
		return c.Redirect("/bookmark")
	}

	c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Bookmark deleted"})
	return c.Redirect("/bookmark")
}
