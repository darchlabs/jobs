package main

import (
	"log"
	"os"

	"github.com/darchlabs/jobs/internal/api/providers"
	providers "github.com/darchlabs/jobs/internal/providers/route"
	"github.com/darchlabs/jobs/internal/storage"
	providerstorage "github.com/darchlabs/jobs/internal/storage/provider"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Get and validate env values
	dbPath := os.Getenv("PATH")
	if dbPath == "" {
		log.Fatal("Invalid DB filepath")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Invalid port")
	}

	// Initialize storage
	s, err := storage.New(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize provider storage
	ps := providerstorage.New(s)

	// Initialize fiber
	api := fiber.New()

	// Configure routers
	providers.Route(api, providers.Context{ProviderStorage: *ps})
}
