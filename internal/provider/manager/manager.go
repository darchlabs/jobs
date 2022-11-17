/* Darch Labs implementation for Keepers */

package providermanager

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/darchlabs/jobs/internal/job"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/robfig/cron"
)

// Interface for each manager provider implementation
type Implementation interface {
	SetupAndRun(job *job.Job) error
}

// Manager interface
type Manager interface {
	Create(job *job.Job) error
}

// Manager stuct
type M struct {
	jobstorage *storage.Job
	client     *ethclient.Client
	privateKey string
}

func NewManager(js *storage.Job, client *ethclient.Client, pk string) *M {
	// get jobs from db
	currentJobs, err := js.List()
	if err != nil {
		// Used log fatal 'cause returning a nil could produce unexpected behaviours
		log.Fatal("cannot get current jobs in the storage")
	}

	m := &M{
		jobstorage: js,
		client:     client,
		privateKey: pk,
	}

	// Iterate jobs and create them for if there were jobs running, sthing failed and is needed to be reloaded
	for _, job := range currentJobs {
		m.Create(job)
	}

	return m
}

// Method for creating a new manager provider
func (m *M) Create(job *job.Job) error {
	var err error

	if job.Type != "cronjob" && job.Type != "synchronizer" {
		return fmt.Errorf("invalid '%s' job type", job.Type)
	}

	// Cronjob based keeper implementation
	if job.Type == "cronjob" {
		cron := cron.New()
		cronjob := NewCronjob(m, cron)

		// TODO(nb): Implement goroutine here?
		err = cronjob.SetupAndRun(job)
		if err != nil {
			return err
		}

		return nil
	}

	// Synchronizer based keeper implementation
	sync := NewSynchronizer()
	err = sync.SetupAndRun(job)
	if err != nil {
		return err
	}

	return nil
}
