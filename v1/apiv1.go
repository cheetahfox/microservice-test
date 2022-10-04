package v1

import "github.com/gofiber/fiber/v2"

func Get(c *fiber.Ctx) error {
	return c.SendString("v1")
}
