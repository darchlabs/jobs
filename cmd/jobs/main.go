package main

import (
	"fmt"
	"log"
	"os"

	jobsapi "github.com/darchlabs/jobs/internal/api/jobs"
	"github.com/darchlabs/jobs/internal/api/providers"
	providermanager "github.com/darchlabs/jobs/internal/provider/manager"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Get and validate env values
	godotenv.Load(".env")

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

	// Initialize provider and job's storage
	js := storage.NewJob(s)
	m := providermanager.NewManager(js)

	// Initialize fiber
	api := fiber.New()

	// Configure routers
<<<<<<< HEAD
	providers.Route(api, providers.Context{ProviderStorage: *ps})
	jobsapi.Route(api, jobsapi.Context{JobStorage: *js, Manager: m})
=======
	providers.Route(api)
	jobsapi.Route(api, jobsapi.Context{JobStorage: *js})
>>>>>>> nb-feat/create-jobs-endpoint

	// Run api
	err = api.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}

}
