package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func checkToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := userRepository.FindTokenById(c.Get("token"))
		if err != nil && token.Id == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		path := strings.Split(c.Path(), "/")[2]
		match := false
		for _, scope := range token.Scope {
			if scope.Name == path {
				match = true
			}
		}
		if !match {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	}
}
