package main

import (
	"fmt"
	"log"
	"os"

	jobsapi "github.com/darchlabs/jobs/internal/api/jobs"
	"github.com/darchlabs/jobs/internal/api/providers"
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

	// Initialize fiber
	api := fiber.New()

	// Configure routers
	providers.Route(api)
	jobsapi.Route(api, jobsapi.Context{JobStorage: *js})

	// Run api
	err = api.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}

}
