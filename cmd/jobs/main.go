package main

import (
	"fmt"
	"log"
	"os"

	jobsapi "github.com/darchlabs/jobs/internal/api/jobs"
	"github.com/darchlabs/jobs/internal/api/providers"
	providermanager "github.com/darchlabs/jobs/internal/provider/manager"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/ethereum/go-ethereum/ethclient"
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

	clientUrl := os.Getenv("URL")
	if clientUrl == "" {
		log.Fatal("Invalid URL")
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if clientUrl == "" {
		log.Fatal("Invalid PRIVATE_KEY")
	}

	// Initialize storage
	s, err := storage.New(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Instance job's storage and client
	js := storage.NewJob(s)
	client, err := ethclient.Dial(clientUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize manager with its params
	m := providermanager.NewManager(js, client, privateKey)

	// Initialize fiber
	api := fiber.New()

	// Configure routers
	providers.Route(api)
	jobsapi.Route(api, jobsapi.Context{JobStorage: *js, Manager: m})

	// Run api
	err = api.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
