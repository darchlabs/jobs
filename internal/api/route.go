package api

import (
	jobstorage "github.com/darchlabs/jobs/internal/storage/job"
	"github.com/gofiber/fiber/v2"
)

type Context struct {
	Storage jobstorage.Storage
}

func Route(app *fiber.App, ctx Context) {
	// TODO(nb): Select and Setup could be both in the same handler?
	app.Post("/api/v1/providers/select/:provider_name", selectProviderHandler(ctx))
	app.Post("/api/v1/providers/setup/:provider_name", setupProviderHandler(ctx))
	app.Post("/api/v1/providers/fund/:provider_name", fundProviderHandler(ctx))
	app.Get("/api/v1/providers/:provder_name?", getProviderHandler(ctx))
}
