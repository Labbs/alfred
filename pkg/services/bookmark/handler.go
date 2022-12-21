package bookmark

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/labbs/alfred/pkg/logger"
)

func (h bookmarkHandler) bookmarkList(c *fiber.Ctx) error {
	d := c.Locals("commonData").(fiber.Map)
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	d["Title"] = "Dashboards"
	d["Error"] = c.Cookies("error-flash")
	d["Success"] = c.Cookies("success-flash")
	c.ClearCookie("error-flash", "success-flash")

	tagFilter := c.Query("tag")
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
		d["Tags"] = []Tag{}
	} else {
		d["Tags"] = tags
	}

	if tagFilter != "" {
		var filteredBookmarks []Bookmark
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
		var searchedBookmarks []Bookmark
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

func (h bookmarkHandler) cleanUnusedTags(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
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

func (h bookmarkHandler) createBookmark(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	c.ClearCookie("Error", "Success")

	var bookmark Bookmark

	name, description, favicon := scraping(c.FormValue("url"))
	if bookmark.Name == "" {
		bookmark.Name = name
	}

	bookmark.Description = description
	bookmark.Icon = favicon

	bookmark.Id = utils.UUIDv4()
	bookmark.UserId = store.Get("user_id").(string)
	bookmark.Url = c.FormValue("url")

	var tags []string
	if c.FormValue("tags") != "" {
		tags = strings.Split(c.FormValue("tags"), ",")
	} else {
		tags = []string{}
	}

	for _, t := range tags {
		tag, err := h.bookmark.GetTagByName(store.Get("user_id").(string), t)
		if err != nil && err.Error.Error() != "record not found" {
			logger.Logger.Error().Err(err.Error).Str("event", "webui.createBookmark").Msg("failed to get tag")
			c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to create bookmark"})
			return c.Redirect("/bookmark")
		}
		if tag.Id == "" {
			tag.Id = utils.UUIDv4()
			tag.Name = t
			tag.UserId = store.Get("user_id").(string)
			err := h.bookmark.CreateTag(tag)
			if err != nil {
				logger.Logger.Error().Err(err.Error).Str("event", "webui.createBookmark").Msg("failed to create tag")
				c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to create bookmark"})
				return c.Redirect("/bookmark")
			}
		}
		bookmark.Tags = append(bookmark.Tags, &tag)
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

func (h bookmarkHandler) createBulkBookmark(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	c.ClearCookie("Error", "Success")

	urls := regexp.MustCompile("\r?\n").Split(c.FormValue("urls"), -1)
	tags := strings.Split(c.FormValue("tags"), ",")

	fmt.Println(urls)

	for _, url := range urls {
		fmt.Println(url)
		var bookmark Bookmark
		name, description, favicon := scraping(url)
		if bookmark.Name == "" {
			bookmark.Name = name
		}

		bookmark.Description = description
		bookmark.Icon = favicon

		bookmark.Id = utils.UUIDv4()
		bookmark.UserId = store.Get("user_id").(string)
		bookmark.Url = url

		for _, t := range tags {
			tag, err := h.bookmark.GetTagByName(store.Get("user_id").(string), t)
			if err != nil && err.Error.Error() != "record not found" {
				logger.Logger.Error().Err(err.Error).Str("event", "webui.createBookmark").Msg("failed to get tag")
				c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to create bookmark"})
				return c.Redirect("/bookmark")
			}
			if tag.Id == "" {
				tag.Id = utils.UUIDv4()
				tag.Name = t
				tag.UserId = store.Get("user_id").(string)
				err := h.bookmark.CreateTag(tag)
				if err != nil {
					logger.Logger.Error().Err(err.Error).Str("event", "webui.createBookmark").Msg("failed to create tag")
					c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to create bookmark"})
					return c.Redirect("/bookmark")
				}
			}
			bookmark.Tags = append(bookmark.Tags, &tag)
		}

		err := h.bookmark.CreateBookmark(bookmark)
		if err != nil {
			logger.Logger.Error().Err(err.Error).Str("event", "webui.createBulkBookmark").Msg("failed to create bookmark")
			c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to create bookmark"})
			return c.Redirect("/bookmark")
		} else {
			logger.Logger.Info().Str("event", "webui.createBulkBookmark").Msg("bookmark created")
			c.Cookie(&fiber.Cookie{Name: "success-flash", Value: "Bookmark created successfully"})
		}
	}
	return c.Redirect("/bookmark")
}

func (h bookmarkHandler) deleteBookmark(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
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

func (h bookmarkHandler) editBookmark(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
	c.ClearCookie("Error", "Success")

	var bookmark Bookmark
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

	var tags []*Tag
	_tags := strings.Split(c.FormValue("tags_list"), ",")
	for _, t := range _tags {
		tag, err := h.bookmark.GetTagByName(store.Get("user_id").(string), t)
		if err != nil && err.Error.Error() != "record not found" {
			logger.Logger.Error().Err(err.Error).Str("event", "webui.editBookmark").Msg("failed to get tag")
			c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to edit bookmark"})
			return c.Redirect("/bookmark")
		}
		if tag.Id == "" {
			tag.Id = utils.UUIDv4()
			tag.Name = t
			tag.UserId = store.Get("user_id").(string)
			err := h.bookmark.CreateTag(tag)
			if err != nil {
				logger.Logger.Error().Err(err.Error).Str("event", "webui.editBookmark").Msg("failed to create tag")
				c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to edit bookmark"})
				return c.Redirect("/bookmark")
			}
		}
		tags = append(tags, &tag)
	}

	for _, t := range _b.Tags {
		exist := false
		for _, _t := range tags {
			if t.Name == _t.Name {
				exist = true
				break
			}
		}
		if !exist {
			err := h.bookmark.DeleteUnusedTag(_b, *t)
			if err != nil {
				logger.Logger.Error().Err(err.Error).Str("event", "webui.editBookmark").Msg("failed to delete unused tag")
				c.Cookie(&fiber.Cookie{Name: "error-flash", Value: "Failed to edit bookmark"})
				return c.Redirect("/bookmark")
			}
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

func (h bookmarkHandler) importBookmarksFile(c *fiber.Ctx) error {
	store, _ := c.Locals("sessions").(*session.Store).Get(c)
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
			var bookmark Bookmark
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
