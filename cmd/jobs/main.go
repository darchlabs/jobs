package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	jobsapi "github.com/darchlabs/jobs/internal/api/jobs"
	"github.com/darchlabs/jobs/internal/api/providers"
	"github.com/darchlabs/jobs/internal/config"
	providermanager "github.com/darchlabs/jobs/internal/provider/manager"
	"github.com/darchlabs/jobs/internal/storage"
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

	// Initialize manager with its params
	m := providermanager.NewManager(js)

	// Initialize fiber
	api := fiber.New()

	// Configure routers
	providers.Route(api)
	jobsapi.Route(api, jobsapi.Context{JobStorage: js, Manager: m})

	// Start already created jobs
	go m.StartCurrentJobs()

	// Run api
	err = api.Listen(fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		log.Fatal(err)
	}

	// Manage shutdowns correctly
	quit := make(chan struct{})
	listenInterrupt(quit)
	<-quit
	gracefullShutdown(m)
}

// listenInterrupt method used to listen SIGTERM OS Signal
func listenInterrupt(quit chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-c
		fmt.Println("Signal received", s.String())
		quit <- struct{}{}
	}()
}

// gracefullShutdown method used to close all synchronizer processes
func gracefullShutdown(m *providermanager.M) {
	// stop all cronjob tickers
	for id := range m.CronMap {
		m.Stop(id)
	}

	// close database connection
	err := m.Jobstorage.Stop()
	if err != nil {
		log.Println(err)
	}
}
