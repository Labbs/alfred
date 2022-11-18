package bookmark

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h bookmarkHandler) bulkEditBookmarks(c *fiber.Ctx) error {
	// _, _ := common.CommonData(h.sessions, c)
	c.ClearCookie("Error", "Success")

	fmt.Println(string(c.Body()))
	return c.Redirect("/bookmark")
}
