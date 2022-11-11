/* Darch Labs implementation for Keepers */

package providermanager

import (
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
		err = m.createCronjob(j)
		return err
	}

	err = m.createSynchronizer(j)
	return err
}

func (m *M) createCronjob(j *job.Job) error {
	err := m.cron.AddFunc(j.Cronjob, func() {
		// Set execute as true for if the job doesn't need to check any method, only needs to work by cronjob time
		execute := true

		// Get blockchain id
		fmt.Println("Getting network..")
		chainId := getChainId(j.Network)
		if chainId == int64(0) {
			err := fmt.Errorf("invalid chain id for %s network", j.Network)
			m.checkAndStop(err, j)
		}

		// Get signer for then execute the tx and evaluate it
		fmt.Println("Getting signer...")
		signer, err := getSigner(m.privateKey, *m.client, chainId, nil, nil)
		m.checkAndStop(err, j)
		fmt.Println("Signer ready!")

		// TODO(nb): j.Abi should be received in this format
		abiFormatted := "[{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"perform\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"setStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

		fmt.Println("Parsing abi...")
		parsedAbi, err := abi.JSON(strings.NewReader(abiFormatted))
		m.checkAndStop(err, j)
		fmt.Println("Abi parsed!")

		fmt.Println("Getting contract...")
		contract := GetContract(j.Address, parsedAbi, m.client)
		m.checkAndStop(err, j)
		fmt.Println("Contract ready!")

		// Check the action method exists
		actionMethod := parsedAbi.Methods[j.ActionMethod].String()
		if actionMethod == "" {
			err = fmt.Errorf("there is no %s method inside the contract abi", actionMethod)
			m.checkAndStop(err, j)
		}

		// check if j.CheckMethod is nil. If it is, execute action method directly
		if j.CheckMethod != nil {
			checkMethod := parsedAbi.Methods[*j.CheckMethod].String()
			if checkMethod == "" {
				err = fmt.Errorf("there is no %s method inside the contract abi", checkMethod)
				m.checkAndStop(err, j)
			}

			fmt.Println("Checking method...")
			res, err := Call(contract, m.client, j.Address, *j.CheckMethod, &bind.CallOpts{})
			m.checkAndStop(err, j)
			fmt.Println("Check method response: ", *res)

			execute = *res
		}

		// Execute action method and see return if it is needed
		if execute && err == nil {
			fmt.Println("Performing tx...")
			tx, err := Perform(contract, m.client, j.Address, j.ActionMethod, &bind.TransactOpts{
				From:     signer.From,
				Signer:   signer.Signer,
				GasLimit: signer.GasLimit,
			}, nil)
			m.checkAndStop(err, j)
			fmt.Println("Tx performed!: ", tx.Hash())
		}
	})

	if err != nil {
		return err
	}

	fmt.Println("Starting cron ...")
	m.cron.Start()
	fmt.Println("Cron started!")

	return nil
}

// Method for stopping cronjob when an error is occurred
func (m *M) checkAndStop(err error, j *job.Job) {
	if err != nil {
		fmt.Println("err: ", err)

		fmt.Println("Updating job...")
		j.Status = fmt.Sprintf("error: %v", err)
		m.jobstorage.Update(j)
		fmt.Println("Job updated!")

		// For stopping cronjob only on that time
		// m.cron.ErrorLog.Fatalf("Cron log err: %v", err)

		// For stopping cron on the following times too, it'll need a restart
		m.cron.Stop()
		fmt.Println("Cron stopped!")
	}
}

// TODO(nb): V2 create syncrhronizer code for listening to an event
// Implementation when a listener over events is needed (synchronizer) --> Jobs V2
func (m *M) createSynchronizer(j *job.Job) error {
	return nil
}
