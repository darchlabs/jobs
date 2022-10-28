package providers

import (
	providerstorage "github.com/darchlabs/jobs/internal/storage/provider"
	"github.com/gofiber/fiber/v2"
)

type Context struct {
	ProviderStorage providerstorage.PS
}

func Route(app *fiber.App, ctx Context) {
	app.Get("/api/v1/jobs/providers", listProvidersHandler(ctx))
	// app.Get("/api/v1/jobs", listJobsHandler(ctx))
	// app.Post("/api/v1/jobs", createProviderHandler(ctx))
}
