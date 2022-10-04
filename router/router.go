package router

import (
	"github.com/cheetahfox/microservice-test/health"
	v1 "github.com/cheetahfox/microservice-test/v1"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Basic Test")
	})
	app.Get("/v1/", v1.Get)
	app.Get("/healthz", health.GetHealthz)
	app.Get("/readyz", health.GetReadyz)

}
