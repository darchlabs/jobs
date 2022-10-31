package jobsapi

import (
	"github.com/darchlabs/jobs/internal/api"
	jobstorage "github.com/darchlabs/jobs/internal/storage/job"
	"github.com/gofiber/fiber/v2"
)

type Context struct {
	JobStorage jobstorage.JS
}

// TODO(nb): Should I put this handler definitions in `api.go` file? Because they're repeated
// in both router files. The same question goes for `handlerFunc` method

// handler's response
type HandlerRes struct {
	Payload    interface{}
	HttpStatus int
	err        error
}

// Define handler method
type handler func(ctx Context) *HandlerRes

func Route(app *fiber.App, ctx Context) {
	listJobsHandler := NewListJobsHandler(ctx.JobStorage)
	app.Get("/api/v1/jobs", HandleFunc(listJobsHandler.Invoke, ctx))
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
