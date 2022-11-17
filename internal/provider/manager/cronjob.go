package providermanager

import (
	"fmt"
	"strings"

	"github.com/darchlabs/jobs/internal/job"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/robfig/cron"
)

type Cronjob struct {
	cron       *cron.Cron
	jobstorage *storage.Job
	client     *ethclient.Client
	pk         string
}

func NewCronjob(ctx *M, cron *cron.Cron) *Cronjob {
	return &Cronjob{
		cron:       cron,
		jobstorage: ctx.jobstorage,
		client:     ctx.client,
		pk:         ctx.privateKey,
	}
}

func (cj *Cronjob) SetupAndRun(job *job.Job) error {
	// Setup cronjob func
	err := cj.cron.AddFunc(job.Cronjob, func() {
		// Set execute as true for if the job doesn't need to check any method, only needs to work by cronjob time
		execute := true

		// Get blockchain id
		fmt.Println("Getting network..")
		chainId := getChainId(job.Network)
		if chainId == int64(0) {
			err := fmt.Errorf("invalid chain id for %s network", job.Network)
			cj.checkAndStop(err, job)
		}

		// Get signer for then execute the tx and evaluate it
		fmt.Println("Getting signer...")
		signer, err := getSigner(cj.pk, *cj.client, chainId, nil, nil)
		cj.checkAndStop(err, job)
		fmt.Println("Signer ready!")

		/// @dev: The abi format must be JSON escaped in order to work correctly
		fmt.Println("Parsing abi...")
		parsedAbi, err := abi.JSON(strings.NewReader(job.Abi))
		cj.checkAndStop(err, job)
		fmt.Println("Abi parsed!")

		fmt.Println("Getting contract...")
		contract := GetContract(job.Address, parsedAbi, cj.client)
		cj.checkAndStop(err, job)
		fmt.Println("Contract ready!")

		// Check the action method exists
		actionMethod := parsedAbi.Methods[job.ActionMethod].String()
		if actionMethod == "" {
			err = fmt.Errorf("there is no %s method inside the contract abi", actionMethod)
			cj.checkAndStop(err, job)
		}

		// check if j.CheckMethod is nil. If it is, execute action method directly
		if job.CheckMethod != nil {
			checkMethod := parsedAbi.Methods[*job.CheckMethod].String()
			if checkMethod == "" {
				err = fmt.Errorf("there is no %s method inside the contract abi", checkMethod)
				cj.checkAndStop(err, job)
			}

			fmt.Println("Checking method...")
			res, err := Call(contract, cj.client, job.Address, *job.CheckMethod, &bind.CallOpts{})
			cj.checkAndStop(err, job)
			fmt.Println("Check method response: ", *res)

			execute = *res
		}

		// Execute action method and see return if it is needed
		if execute && err == nil {
			fmt.Println("Performing tx...")
			tx, err := Perform(contract, cj.client, job.Address, job.ActionMethod, &bind.TransactOpts{
				From:     signer.From,
				Signer:   signer.Signer,
				GasLimit: signer.GasLimit,
			}, nil)
			cj.checkAndStop(err, job)
			fmt.Println("Tx performed!: ", tx.Hash())
		}
	})

	// Check there is no error in the cronjob func
	if err != nil {
		return err
	}

	// Execute cronjob func
	cj.cron.Start()
	fmt.Println("Cron started!")

	return nil
}

// Method for stopping cronjob when an error is occurred
func (cj *Cronjob) checkAndStop(err error, j *job.Job) {
	if err != nil {
		fmt.Println("err: ", err)

		// Update the status of the job on the storage
		fmt.Println("Updating job...")
		j.Status = fmt.Sprintf("error: %v", err)
		cj.jobstorage.Update(j)
		fmt.Println("Job updated!")

		// It stops the cronjob and then it'll need to be re-executed to start again
		cj.cron.Stop()
		fmt.Println("Cron stopped!")
	}
}
