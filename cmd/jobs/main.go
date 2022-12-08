package main

import (
	"fmt"
	"log"

	jobsapi "github.com/darchlabs/jobs/internal/api/jobs"
	"github.com/darchlabs/jobs/internal/api/providers"
	"github.com/darchlabs/jobs/internal/config"
	providermanager "github.com/darchlabs/jobs/internal/provider/manager"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofiber/fiber/v2"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	log.Println("Starting jobs")

	// read env values
	var conf config.Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatal("invalid env values: ", err)
	}

	// Initialize storage
	s, err := storage.New(conf.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	// Instance job's storage and client
	js := storage.NewJob(s)
	client, err := ethclient.Dial(conf.NodeURL)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize manager with its params
	m := providermanager.NewManager(js, client, conf.PrivateKey)

	// Initialize fiber
	api := fiber.New()

	// Configure routers
	providers.Route(api)
	jobsapi.Route(api, jobsapi.Context{JobStorage: js, Manager: m})

	// Run api
	err = api.Listen(fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		log.Fatal(err)
	}
}
