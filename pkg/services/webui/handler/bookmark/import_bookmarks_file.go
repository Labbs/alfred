package bookmark

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/alfred/pkg/common"
	"github.com/labbs/alfred/pkg/logger"
	b "github.com/labbs/alfred/pkg/services/bookmark"
)

func (h bookmarkHandler) importBookmarksFile(c *fiber.Ctx) error {
	_, store := common.CommonData(h.sessions, c)
	c.ClearCookie("Error", "Success")

	form, errf := c.MultipartForm()
	if errf != nil {
		logger.Logger.Error().Err(errf).Str("event", "bookmark.import_bookmarks_file").Msg("could_not_get_bookmarks_file")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get bookmarks file"})
		return c.Redirect("/user/profile")
	}

	e := c.SaveFileToStorage(form.File["bookmark"][0], "bookmarks", temporaryStore)
	if e != nil {
		logger.Logger.Error().Err(e).Str("event", "bookmark.import_bookmarks_file").Msg("could_not_save_bookmarks_file")
		c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not save bookmarks file"})
		return c.Redirect("/user/profile")
	}

	go func() {
		file, errS := temporaryStore.Get("bookmarks")
		if errS != nil {
			logger.Logger.Error().Err(errS).Str("event", "bookmark.import_bookmarks_file").Msg("could_not_get_bookmarks_file_from_memory_storage")
			c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not get bookmarks file from memory storage"})
		}

		links := []string{}

		for _, line := range strings.Split(string(file), "\n") {
			if strings.Contains(line, "HREF=") {
				links = append(links, strings.ToLower(strings.Split(line, "<DT>")[1]))
			}
		}

		for _, link := range links {
			var bookmark b.Bookmark
			bookmark.Id = utils.UUIDv4()
			bookmark.UserId = store.Get("user_id").(string)

			reHref := regexp.MustCompile(".*href=\"(.*?)\" .*")
			bookmark.Url = reHref.FindStringSubmatch(link)[1]

			reName := regexp.MustCompile(".*>(.*)<.*")
			bookmark.Name = reName.FindStringSubmatch(link)[1]

			if strings.Contains(link, "icon") {
				reIcon := regexp.MustCompile(".*icon=\"(.*?)\".*")
				bookmark.Icon = reIcon.FindStringSubmatch(link)[1]
			}

			err := h.bookmark.CreateBookmark(bookmark)
			if err != nil {
				logger.Logger.Error().Err(err.Error).Str("event", "bookmark.import_bookmarks_file").Msg("could_not_create_bookmark")
				c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Could not create bookmark"})
			}
		}

		logger.Logger.Info().Str("event", "bookmark.import_bookmarks_file").Msg("bookmarks_imported")
		c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Bookmarks imported successfully"})
	}()

	logger.Logger.Info().Str("event", "bookmark.import_bookmarks_file").Msg("bookmarks_import_in_progress")
	c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Bookmarks import in progress"})

	return c.Redirect("/bookmark")
}
