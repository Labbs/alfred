package bookmark

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
	b "github.com/labbs/alfred/pkg/services/bookmark"
)

func (h bookmarkHandler) createBookmark(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
	c.ClearCookie("Error", "Success")

	var bookmark b.Bookmark

	name, description, favicon := scraping(c.FormValue("url"))
	if bookmark.Name == "" {
		bookmark.Name = name
	}

	bookmark.Description = description
	bookmark.Icon = favicon

	bookmark.Id = utils.UUIDv4()
	bookmark.UserId = store.Get("user_id").(string)
	bookmark.Url = c.FormValue("url")

	tags := strings.Split(c.FormValue("tags"), ",")
	if len(tags) != 0 {
		for _, t := range tags {
			bookmark.Tags = append(bookmark.Tags, b.Tag{
				Id:     utils.UUIDv4(),
				Name:   t,
				UserId: store.Get("user_id").(string)})
		}
	}

	err := h.bookmark.CreateBookmark(bookmark)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "webui.createBookmark").Msg("failed to create bookmark")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to create bookmark"})
		return c.Redirect("/bookmark")
	} else {
		logger.Logger.Info().Str("event", "webui.createBookmark").Msg("bookmark created")
		c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Bookmark created successfully"})
		return c.Redirect("/bookmark")
	}
}
