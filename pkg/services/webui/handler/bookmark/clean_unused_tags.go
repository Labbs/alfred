package bookmark

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
)

func (h bookmarkHandler) cleanUnusedTags(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
	c.ClearCookie("Error", "Success")

	tags, err := h.bookmark.GetTags(store.Get("user_id").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "webui.cleanUnusedTags").Msg("failed to get tags")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to clean unused tags"})
		return c.Redirect("/bookmark")
	}

	for _, t := range tags {
		if len(t.Bookmarks) == 0 {
			err := h.bookmark.DeleteTag(t.Id, store.Get("user_id").(string))
			if err != nil {
				logger.Logger.Error().Err(err.Error).Str("event", "webui.cleanUnusedTags").Msg("failed to delete tag")
				c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to clean unused tags"})
				return c.Redirect("/bookmark")
			}
		}
	}

	c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Unused tags cleaned"})
	return c.Redirect("/bookmark")
}
