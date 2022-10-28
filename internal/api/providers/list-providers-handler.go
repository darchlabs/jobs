package providers

import (
	"github.com/darchlabs/jobs/internal/api"
	"github.com/gofiber/fiber/v2"
)

func listProvidersHandler(ctx Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		// Get elements from db
		data, err := ctx.ProviderStorage.List()
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(api.Response{
				Error: err.Error(),
			})
		}

		// Prepare response
		return c.Status(fiber.StatusOK).JSON(api.Response{
			Data: data,
		})
	}

}
