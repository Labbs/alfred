package server

import (
	"embed"
	"net/http"

	"github.com/gofiber/template/html"
)

var (
	//go:embed templates/*
	templatesfs embed.FS

	//go:embed static/*
	embedDirStatic embed.FS
)

func engineInit() *html.Engine {
	engine := html.NewFileSystem(http.FS(templatesfs), ".html")
	engine.AddFunc("add", add)
	engine.AddFunc("currentVersion", currentVersion)
	engine.AddFunc("truncate", truncate)
	engine.AddFunc("truncateByWord", truncateByWord)
	engine.AddFunc("joinTags", joinTags)
	engine.AddFunc("toJson", toJson)
	engine.AddFunc("widgetConfigTransform", widgetConfigTransform)
	return engine
}
