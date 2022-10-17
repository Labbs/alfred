package logger

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitLogger(version string) {
	host, _ := os.Hostname()
	Logger = zerolog.
		New(os.Stderr).
		With().
		Caller().
		Timestamp().
		Str("host", host).
		Str("version", version).
		Logger()
}

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		timeStart := time.Now()
		c.Next()
		Logger.Info().
			Int("status", c.Response().StatusCode()).
			Dur("duration", time.Since(timeStart)).
			Str("method", string(c.Request().Header.Method())).
			Str("remote_addr", c.IP()).
			Str("path", c.Path()).
			Str("user_agent", c.Get("User-Agent")).
			Msg("")
		return nil
	}
}
