/*
Implement health checks for kubernetes Liveness and Readiness.
*/
package health

import (
	"github.com/gofiber/fiber/v2"
)

var ApiReady bool

func init() {
	ApiReady = false
}

func GetHealthz(c *fiber.Ctx) error {
	// return &fiber.Error{}
	return c.SendStatus(200)
}

func GetReadyz(c *fiber.Ctx) error {
	if !ApiReady {
		return c.SendStatus(503)
	}
	return c.SendStatus(200)
}
