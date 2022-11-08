/* Darch Labs implementation for Keepers */

package providermanager

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/darchlabs/jobs/internal/job"
	"github.com/darchlabs/jobs/internal/storage"
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
	fmt.Println("Entered")
	err := m.cron.AddFunc(j.Cronjob, func() {
		execute := true

		fmt.Println("I'm here!")

		// Get blockchain id
		chainId := getChainId(j.Network)
		if chainId == int64(0) {
			err := fmt.Errorf("invalid chain id for %s network", j.Network)
			m.checkAndStop(err, j)
		}

		// Get signer for then execute the tx and evaluate it
		signer, err := getSigner(m.privateKey, *m.client, chainId, nil, nil)
		m.checkAndStop(err, j)

		// Parse address
		address := common.HexToAddress(j.Address)

		// Get an instance of the smart contract
		contract, err := NewToken(address, j.Abi, m.client)
		m.checkAndStop(err, j)

		// check if j.CheckMethod is nil. If it is nil execute action methdo directly.
		if j.CheckMethod != nil {
			// Create results interface array for appending the result of the sc view function on it
			results := make([]interface{}, 0)

			// Call the view function and get its boolean response
			err = contract.contractCaller.Call(&bind.CallOpts{}, &results, *j.CheckMethod, nil)
			m.checkAndStop(err, j)

			execute = results[0] == true
		}

		// Execute action method and see return
		if execute && err == nil {
			fmt.Print("1")
			tx, err := contract.contractTransactor.Transact(&bind.TransactOpts{From: signer.From}, j.ActionMethod, nil)

			m.checkAndStop(err, j)
			fmt.Println("tx: ", tx)
		}

		// TODO(nb): add state or chan for running different process
	})

	if err != nil {
		return err
	}

	m.cron.Start()
	return nil
}

// Method for stopping cronjob when an error is occurred
func (m *M) checkAndStop(err error, j *job.Job) {
	if err != nil {
		j.Status = "error"
		m.jobstorage.Update(j)
		m.cron.Stop()
	}
}

// TODO(nb): V2 create syncrhronizer code for listening to an event
// Implementation when a listener over events is needed (synchronizer) --> Jobs V2
func (m *M) createSynchronizer(j *job.Job) error {
	return nil
}
