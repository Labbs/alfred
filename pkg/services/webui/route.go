package webui

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"github.com/labbs/fiber-storage/gorm"

	"github.com/labbs/alfred/pkg/database"
	"github.com/labbs/alfred/pkg/services/webui/handler/auth"
	book "github.com/labbs/alfred/pkg/services/webui/handler/bookmark"
	dash "github.com/labbs/alfred/pkg/services/webui/handler/dashboard"
	user "github.com/labbs/alfred/pkg/services/webui/handler/user"
)

var (
	sessions *session.Store

	//go:embed templates/*
	templatesfs embed.FS

	//go:embed static/*
	embedDirStatic embed.FS
)

func EngineInit() *html.Engine {
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

func InitRoute(r *fiber.App) {
	storage := gorm.New(gorm.Config{
		DB:    database.GetDbConnection().DB,
		Table: "session",
	})
	sessions = session.New(session.Config{
		Storage: storage,
	})

	r.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirStatic),
		PathPrefix: "static",
		Browse:     true,
	}))

	g := r.Group("/")
	g.Use(checkSession())
	auth.InitRoute(g, sessions)
	dash.InitRoute(g, sessions)
	book.InitRoute(g, sessions)
	user.InitRoute(g, sessions)
}
