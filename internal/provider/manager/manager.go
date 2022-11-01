package providermanager

import (
	"fmt"
	"log"

	"github.com/darchlabs/jobs/internal/job"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/robfig/cron"
)

type Manager interface {
	Create(j *job.Job) error
}

type M struct {
	jobs []*job.Job
	cron *cron.Cron
}

func NewManager(js *storage.Job) *M {
	// get jobs from db
	currentJobs, err := js.List()
	if err != nil {
		// Used log fatal 'cause returning a nil could produce unexpected behaviours
		log.Fatal("cannot get current jobs in the storage")
	}

	m := &M{
		jobs: make([]*job.Job, 0),
		cron: cron.New(),
	}

	// Iterate jobs and create them for if there were jobs running, sthing failed and is needed to be reloaded
	for _, job := range currentJobs {
		m.Create(job)
	}

	return m
}

func (m *M) Create(j *job.Job) error {
	m.jobs = append(m.jobs, j)
	var err error

	if j.Type != "cronjob" && j.Type != "synchronizer" {
		return fmt.Errorf("invalid '%s' job type", j.Type)
	}

	if j.Type == "cronjob" {
		fmt.Println("Yes")
		err = m.createCronjob(j)
		return err
	}

	err = m.createSynchronizer(j)
	return err
}

func (m *M) createCronjob(j *job.Job) error {
	fmt.Println("Entered")
	m.cron.AddFunc(j.Cronjob, func() {
		fmt.Println("I'm here!")

		// Get provider and evaluate it

		// Get instance of the contract w address (and abi)

		// check if j.CheckMethod is nil. If it is nil execute action methdo directly.
		// If not, execute action method and see return

		// TODO(nb): add state or chan for running different process

		// TODO(nb): if the status changes (e.g error) the BD needs to be updated. Decide where & how
	})
	m.cron.Start()

	return nil
}

// Implementation when a listener over events is needed (synchronizer) --> V2
func (m *M) createSynchronizer(j *job.Job) error {
	// TODO(nb): create syncrhronizer code for listening to an event
	return nil
}
