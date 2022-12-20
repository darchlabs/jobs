package providers

import (
	"github.com/darchlabs/jobs/internal/api"
	"github.com/gofiber/fiber/v2"
)

// Define handler method
type handler func() *api.HandlerRes

func Route(app *fiber.App) {
	listProvidersHandler := NewListProvidersHandler()

	app.Get("/api/v1/jobs/providers", HandleFunc(listProvidersHandler.Invoke))
}

// Func that receives the returns from handlers and creates an http response
func HandleFunc(fn handler) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// set headers
		c.Accepts("application/json")

		// Exec handler func and get its response
		handlerRes := fn()
		payload, statusCode, err := handlerRes.Payload, handlerRes.HttpStatus, handlerRes.Err
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(api.Response{
				Data:  payload,
				Meta:  statusCode,
				Error: err,
			})
		}

		// Prepare response
		res := api.Response{
			Meta: map[string]interface{}{"statusCode": statusCode},
			Data: payload,
		}

		// Return response
		return c.Status(fiber.StatusOK).JSON(res)
	}
}
