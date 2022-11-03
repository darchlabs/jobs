/* Darch Labs implementation for Keepers */

package providermanager

import (
	"fmt"
	"log"

	"github.com/darchlabs/jobs/internal/job"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/robfig/cron"
)

type Manager interface {
	Create(j *job.Job) error
}

type M struct {
	jobstorage *storage.Job
	cron       *cron.Cron
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
		cron:       cron.New(),
		client:     client,
		privateKey: pk,
	}

	// Iterate jobs and create them for if there were jobs running, sthing failed and is needed to be reloaded
	for _, job := range currentJobs {
		m.Create(job)
	}

	return m
}

func (m *M) Create(j *job.Job) error {
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
	var err error
	fmt.Println("Entered")
	err = m.cron.AddFunc(j.Cronjob, func() {
		execute := true

		fmt.Println("I'm here!")

		// Get blockchain id
		chainId := getChainId(j.Network)
		if chainId == int64(0) {
			err = fmt.Errorf("invalid chain id for %s network", j.Network)
			m.stopCronjob(j)
		}

		// Get signer for then execute the tx and evaluate it
		signer, err := getSigner(m.privateKey, *m.client, chainId, nil, nil)
		if err != nil {
			m.stopCronjob(j)
		}

		// TODO(nb): Get instance of the contract w address (and abi)

		// check if j.CheckMethod is nil. If it is nil execute action methdo directly.
		if j.CheckMethod != nil {
			// TODO(nb): Execute smart contract view method
			execute = false
		}

		// Execute action method and see return
		if execute && err == nil {
			fmt.Print("1")
			// TODO(nb): Execute smart contract method
		}

		// TODO(nb): add state or chan for running different process

		if err != nil {
			m.stopCronjob(j)
		}
	})

	if err != nil {
		return err
	}

	m.cron.Start()
	return nil
}

// Method for stopping cronjob when an error is occurred
func (m *M) stopCronjob(j *job.Job) {
	j.Status = "error"
	m.jobstorage.Update(j)
	m.cron.Stop()
}

// Implementation when a listener over events is needed (synchronizer) --> Jobs V2
func (m *M) createSynchronizer(j *job.Job) error {
	// TODO(nb): V2 create syncrhronizer code for listening to an event
	return nil
}
