package bookmark

import (
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/alfred/pkg/logger"
	b "github.com/labbs/alfred/pkg/services/bookmark"
)

func (h bookmarkHander) createBookmark(c *fiber.Ctx) error {
	token, err := h.user.FindTokenById(c.Get("token"))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "api.bookmark.create").Msg("failed to find token")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var bookmark b.Bookmark
	bookmark.UserId = token.UserId
	bookmark.Id = utils.UUIDv4()

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		logger.Logger.Error().Err(err).Str("event", "api.bookmark.create").Msg("failed to parse body")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	bookmark.Name = body["name"].(string)
	bookmark.Url = body["url"].(string)
	bookmark.Icon = body["icon"].(string)

	tags := strings.Split(body["tags"].(string), ",")
	if len(tags) != 0 {
		for _, t := range tags {
			bookmark.Tags = append(bookmark.Tags, b.Tag{
				Id:     utils.UUIDv4(),
				Name:   t,
				UserId: token.UserId})
		}
	}

	err = h.bookmark.CreateBookmark(bookmark)
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "api.bookmark.create").Msg("failed to create bookmark")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "bookmark_created"})
}
