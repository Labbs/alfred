package bookmark

import (
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
	b "github.com/labbs/alfred/pkg/services/bookmark"
)

func (h bookmarkHandler) editBookmark(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
	c.ClearCookie("Error", "Success")

	var bookmark b.Bookmark
	if err := c.BodyParser(&bookmark); err != nil {
		logger.Logger.Error().Err(err).Str("event", "webui.editBookmark").Msg("failed to parse bookmark")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to parse bookmark"})
		return c.Redirect("/bookmark")
	}

	_b, err := h.bookmark.GetBookmarkById(store.Get("user_id").(string), c.Params("id"))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "bookmark.update.get_by_id").Msg("could_not_get_bookmark")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get bookmark"})
		return c.Redirect("/bookmark")
	}

	var tags []b.Tag
	_tags := strings.Split(c.FormValue("tags_list"), ",")
	for _, tag := range _tags {
		exists := false
		for _, t := range _b.Tags {
			toDelete := true
			if t.Name == tag {
				tags = append(tags, t)
				exists = true
				toDelete = false
			}
			if toDelete {
				err := h.bookmark.DeleteTag(t.Id, store.Get("user_id").(string))
				if err != nil {
					logger.Logger.Error().Err(err.Error).Str("event", "bookmark.update.delete_tag").Msg("could_not_delete_tag")
					c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not delete tag"})
					return c.Redirect("/bookmark")
				}
			}
		}
		if !exists {
			tags = append(tags, b.Tag{
				Id:     utils.UUIDv4(),
				Name:   tag,
				UserId: store.Get("user_id").(string)})
		}
	}

	bookmark.Tags = tags
	bookmark.UserId = store.Get("user_id").(string)
	bookmark.Description = c.FormValue("description")
	bookmark.Icon = c.FormValue("icon")
	bookmark.Name = c.FormValue("name")
	bookmark.Id = c.Params("id")

	err = h.bookmark.UpdateBookmark(bookmark)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "webui.editBookmark").Msg("failed to update bookmark")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to update bookmark"})
		return c.Redirect("/bookmark")
	} else {
		logger.Logger.Info().Str("event", "webui.editBookmark").Msg("bookmark updated")
		c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Bookmark updated successfully"})
		return c.Redirect("/bookmark")
	}
}
