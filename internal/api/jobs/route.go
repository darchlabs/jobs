package jobsapi

import (
	jobstorage "github.com/darchlabs/jobs/internal/storage/job"
	"github.com/gofiber/fiber/v2"
)

type Context struct {
	JobStorage jobstorage.JS
}

func Route(app *fiber.App, ctx Context) {
	app.Get("/api/v1/jobs", listJobsHandler(ctx))
	// app.Post("/api/v1/jobs", createProviderHandler(ctx))
}
