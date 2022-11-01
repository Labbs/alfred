package bookmark

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
	b "github.com/labbs/alfred/pkg/services/bookmark"
)

func (h bookmarkHandler) bookmarkList(c *fiber.Ctx) error {
	d, store := common.CommonData(h.sessions, c)
	d["Title"] = "Dashboards"
	d["Error"] = c.Cookies("error-flash")
	d["Success"] = c.Cookies("success-flash")
	c.ClearCookie("error-flash", "success-flash")

	tagFilter := c.Query("tag")
	fmt.Println(c.FormValue("search"))
	search := ""
	if c.Method() == "POST" {
		search = c.FormValue("search")
		d["Search"] = search
	}

	bookmarks, err := h.bookmark.GetAllBookmarks(store.Get("user_id").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "webui.bookmarks").Msg("failed to get all bookmarks")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to get all bookmarks"})
	}

	tags, err := h.bookmark.GetUniqueTags(store.Get("user_id").(string))
	if err != nil {
		logger.Logger.Error().Err(err.Error).Str("event", "webui.bookmarks").Msg("failed to get all tags")
		d["Error"] = "Failed to get all tags"
		d["Tags"] = []b.Tag{}
	} else {
		d["Tags"] = tags
	}

	if tagFilter != "" {
		var filteredBookmarks []b.Bookmark
		for _, bookmark := range bookmarks {
			for _, tag := range bookmark.Tags {
				if tag.Name == tagFilter {
					filteredBookmarks = append(filteredBookmarks, bookmark)
				}
			}
		}
		bookmarks = filteredBookmarks
	}

	if search != "" {
		var searchedBookmarks []b.Bookmark
		for _, bookmark := range bookmarks {
			if strings.Contains(bookmark.Name, search) ||
				strings.Contains(bookmark.Url, search) ||
				strings.Contains(bookmark.Description, search) {
				searchedBookmarks = append(searchedBookmarks, bookmark)
			}
		}

		bookmarks = searchedBookmarks
	}

	d["Bookmarks"] = bookmarks

	return c.Render("templates/bookmark", d, "templates/layouts/main")
}
