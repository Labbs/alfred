package bookmark

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/labbs/alfred/pkg/logger"
)

func scraping(s string) (string, string, string) {
	var title string
	var description string
	var fff []string
	var favicon string

	c := colly.NewCollector(colly.MaxDepth(1))

	c.OnHTML("head title", func(e *colly.HTMLElement) {
		title = e.Text
	})

	c.OnHTML("head meta[property=\"og:description\"]", func(e *colly.HTMLElement) {
		description = e.Attr("content")
	})

	c.OnHTML("head link[rel=\"icon\"]", func(e *colly.HTMLElement) {
		fff = append(fff, e.Attr("href"))
	})

	c.OnHTML("head link[rel=\"shortcut icon\"]", func(e *colly.HTMLElement) {
		fff = append(fff, e.Attr("href"))
	})

	c.OnHTML("head meta[itemprop=\"image\"]", func(e *colly.HTMLElement) {
		fff = append(fff, e.Attr("content"))
	})

	if err := c.Visit(s); err != nil {
		logger.Logger.Error().Err(err).Str("event", "scrapping.visite").Msg(err.Error())
	}

	for _, f := range fff {
		if f != "" && favicon == "" {
			if !strings.Contains(f, "http") {
				u, _ := url.Parse(s)
				favicon = fmt.Sprintf("%s://%s/%s", u.Scheme, u.Hostname(), f)
			}
		}
	}

	return title, description, favicon
}
