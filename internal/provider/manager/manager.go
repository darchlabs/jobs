/* Darch Labs implementation for Keepers */

package providermanager

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/darchlabs/jobs/internal/job"
	"github.com/darchlabs/jobs/internal/provider"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/robfig/cron"
)

// Interface for each manager provider implementation
type Implementation interface {
	SetupAndRun(job *job.Job) error
}

// Manager interface
type Manager interface {
	Setup(job *job.Job) error
	Start(id string)
	StartCurrentJobs()
}

// Manager stuct
type M struct {
	Jobstorage *storage.Job
	client     *ethclient.Client
	privateKey string
	CronMap    map[string]*cron.Cron
}

func NewManager(js *storage.Job, client *ethclient.Client, pk string) *M {
	cronMap := make(map[string]*cron.Cron)
	m := &M{
		Jobstorage: js,
		client:     client,
		privateKey: pk,
		CronMap:    cronMap,
	}

	return m
}

func (m *M) StartCurrentJobs() {
	// get jobs from db
	currentJobs, err := m.Jobstorage.List()
	if err != nil {
		// Used log fatal 'cause returning a nil could produce unexpected behaviours
		log.Fatal("cannot get current jobs in the storage")
	}

	for _, job := range currentJobs {
		err := m.Setup(job)
		if err != nil {
			fmt.Printf("Error while setting up %s job\n", job.ID)
			continue
		}

		if job.Status != provider.StatusRunning {
			continue
		}

		m.Start(job.ID)
	}
}

// Method for creating a new manager provider
func (m *M) Setup(job *job.Job) error {
	if job.Type != "cronjob" && job.Type != "synchronizer" {
		return fmt.Errorf("invalid '%s' job type", job.Type)
	}

	newCron := cron.New()

	// Cronjob based keeper implementation
	if job.Type == "cronjob" {
		cronjob := NewCronjob(m, newCron)
		m.CronMap[job.ID] = newCron

		// Check if the inputs for the cron are right
		cronCTX, err := cronjob.Check(job)
		if err != nil {
			return err
		}

		// Add job to the cron
		stop := make(chan bool)
		err = cronjob.AddJob(job, cronCTX, stop)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *M) Start(id string) {
	c := m.CronMap[id]

	fmt.Println("Starting cron: ", id)
	c.Start()
	fmt.Println("Cron started!")
}

func (m *M) Stop(id string) {
	c := m.CronMap[id]

	fmt.Println("Stopping cron: ", id)
	c.Stop()
	fmt.Println("Cron stopped!")
}
