package providers

import (
	"github.com/darchlabs/jobs/internal/api"
	providerstorage "github.com/darchlabs/jobs/internal/storage/provider"
	"github.com/gofiber/fiber/v2"
)

type Context struct {
	ProviderStorage providerstorage.PS
}

// handler's response
type HandlerRes struct {
	Payload    interface{}
	HttpStatus int
	err        error
}

// Define handler method
type handler func(ctx Context) *HandlerRes

func Route(app *fiber.App, ctx Context) {
	listProvidersHandler := NewListProvidersHandler(ctx.ProviderStorage)

	app.Get("/api/v1/jobs/providers", HandleFunc(listProvidersHandler.Invoke, ctx))
}

// Func that receives the returns from handlers and creates an http response
func HandleFunc(fn handler, ctx Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// set headers
		c.Accepts("application/json")

		// Exec handler func and get its response
		handlerRes := fn(ctx)
		payload, statusCode, err := handlerRes.Payload, handlerRes.HttpStatus, handlerRes.err
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(api.Response{
				Error: err,
				Data:  payload,
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
