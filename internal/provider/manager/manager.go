/* Darch Labs implementation for Keepers */

package providermanager

import (
	"fmt"
	"log"

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
	Stop(id string)
	StartCurrentJobs()
}

// Manager stuct
type M struct {
	Jobstorage *storage.Job
	CronMap    map[string]*cron.Cron
}

func NewManager(js *storage.Job) *M {
	cronMap := make(map[string]*cron.Cron)
	m := &M{
		Jobstorage: js,
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

	// TODO(nb): add gorotuines to each loop iteration for making it fast?
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
	// Save the current cron, for if the new one fails to comeback to it
	currentCron := m.CronMap[job.ID]

	// Create new cron and cronjob instances
	newCron := cron.New()
	cronjob := NewCronjob(m, newCron)

	// Update cron map instance with new cron
	m.CronMap[job.ID] = newCron

	// Check if the inputs for the cron are right
	cronCTX, err := cronjob.Check(job)
	if err != nil {
		fmt.Println("err while checking job: ", err)

		// The cronjob will keep being the currentCron, not the new one
		m.CronMap[job.ID] = currentCron

		// Get job in DB for knowing if it's already created
		job, dbErr := m.Jobstorage.GetById(job.ID)
		if dbErr != nil {
			fmt.Println("dbErr: ", dbErr)
			// If it is not created, it won't update
			return err
		}

		// The error is used as log
		log := err.Error()
		job.Logs = &log

		// It updates the log field to the job in the db
		_, updateErr := m.Jobstorage.Update(job)
		if updateErr != nil {
			fmt.Println("updateErr: ", updateErr)
		}

		return err
	}

	// Add job to the cron
	stop := make(chan bool)
	err = cronjob.AddJob(job, cronCTX, stop)
	if err != nil {
		m.CronMap[job.ID] = currentCron
		return err
	}

	return nil
}

func (m *M) Start(id string) {
	// Get the cron instance of that job id
	c := m.CronMap[id]

	fmt.Println("Starting cron: ", id)
	// It'll wait the cronjob period to pass before starting the 1st job
	c.Start()
	fmt.Println("Cron started!")
}

func (m *M) Stop(id string) {
	// Get the cron instance of that job id
	c := m.CronMap[id]

	fmt.Println("Stopping cron: ", id)
	c.Stop()
}
