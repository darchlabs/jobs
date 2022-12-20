package providermanager

import (
	"context"
	"fmt"
	"strings"

	"github.com/darchlabs/jobs/internal/job"
	"github.com/darchlabs/jobs/internal/provider"
	sc "github.com/darchlabs/jobs/internal/provider/smart-contracts"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/robfig/cron"
)

type Cronjob struct {
	cron       *cron.Cron
	jobstorage *storage.Job
}

func NewCronjob(ctx *M, cron *cron.Cron) *Cronjob {
	return &Cronjob{
		cron:       cron,
		jobstorage: ctx.Jobstorage,
	}
}

// ctx for running an already configured cronjob cron
type cronCTX struct {
	signer   *bind.TransactOpts
	abi      abi.ABI
	contract *bind.BoundContract
	client   *ethclient.Client
}

// Func for asserting that the data is ok, and returning the needed context for the cronjob
func (cj *Cronjob) Check(job *job.Job) (*cronCTX, error) {
	// Declare variables for then using or returning them after adding the cron func
	var err error
	var signer *bind.TransactOpts
	var parsedAbi abi.ABI
	var contract *bind.BoundContract

	// Setup cronjob func
	fmt.Println("Getting network..")
	chainId := getChainId(job.Network)
	if chainId == int64(0) {
		return nil, fmt.Errorf("invalid chain id for %s network", job.Network)
	}
	fmt.Println("Network obtained!")

	client, err := ethclient.Dial(job.NodeURL)
	if err != nil {
		return nil, err
	}

	fmt.Println("Getting signer...")
	signer, err = sc.GetSigner(job.Privatekey, *client, chainId, nil, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Signer ready!")

	/// @dev: The abi format must be JSON escaped in order to work correctly
	fmt.Println("Parsing abi...")
	parsedAbi, err = abi.JSON(strings.NewReader(job.Abi))
	if err != nil {
		return nil, err
	}
	fmt.Println("Abi parsed!")

	// Check that the address indeed is deployed on the network
	contractCode, err := client.PendingCodeAt(context.Background(), common.HexToAddress(job.Address))
	if err != nil {
		return nil, err
	}

	if len(contractCode) == 0 {
		return nil, fmt.Errorf("%s", "the contract address doesn't exist")
	}

	fmt.Println("Getting contract...")
	contract = sc.GetContract(job.Address, parsedAbi, client)

	if contract == nil {
		return nil, fmt.Errorf("there is no contract under the %s address and abi on the %s network", job.Address, job.Network)
	}
	fmt.Println("Contract ready!")

	// Check the action method exists
	fmt.Println("Getting actionMethod...")
	actionMethod := parsedAbi.Methods[job.ActionMethod].String()
	if actionMethod == "" {
		return nil, fmt.Errorf("there is no %s method inside the contract abi", actionMethod)
	}
	fmt.Println("actionMethod is OK!")

	// check if j.CheckMethod is nil. If it is, perform action method directly
	fmt.Println("Getting checkMethod...")
	if job.CheckMethod != nil {
		checkMethod := parsedAbi.Methods[*job.CheckMethod].String()
		if checkMethod == "" {
			return nil, fmt.Errorf("there is no %s method inside the contract abi", checkMethod)
		}
	}
	fmt.Println("checkMethod is OK!")

	return &cronCTX{
		signer:   signer,
		abi:      parsedAbi,
		contract: contract,
		client:   client,
	}, err
}

func (cj *Cronjob) AddJob(job *job.Job, ctx *cronCTX, stop chan bool) error {
	// define log for making the errors more explicit
	var log string
	var err error
	var errCounter uint8

	// Set max times that the cron func can fail before being stopped
	maxErrorsLimit := uint8(5)

	// define variable for the cronjob to know if it must perform or not the job
	perform := true

	// Create the cronjob func
	err = cj.cron.AddFunc(job.Cronjob, func() {
		/* actionMethod check */

		// TODO(nb): handle when the client dies, update error and log
		// if j.CheckMethod is nil, avoid this check
		if job.CheckMethod != nil {
			// Check if the response of the smart contract view function for the cronjob to know if it must perform actionMethod or not
			fmt.Println("Checking method...")
			res, err := sc.Call(ctx.contract, ctx.client, job.Address, *job.CheckMethod, &bind.CallOpts{})
			if err != nil {
				errCounter += 1
				fmt.Println("errCounter: ", errCounter)

				log = fmt.Sprintf("Error while trying to call checkMethod: %v", err)

				if errCounter > maxErrorsLimit {
					fmt.Println("here")
					stopLog := fmt.Sprintf("Failed 3 times so stopped job. Last error: %s", log)
					updateJob(cj.jobstorage, job, stopLog)

					stop <- true
					return
				}

				// Check if the job still doesn't exist, it doesn't need an update
				_, err := cj.jobstorage.GetById(job.ID)
				if err != nil {
					return
				}

				updateJob(cj.jobstorage, job, log)
				return
			}
			fmt.Println("Check method response: ", *res)

			perform = *res
		}

		/* actionMethod perform tx */

		// perform action method and see return if it is needed
		if perform && err == nil {
			fmt.Println("Performing tx...")
			tx, err := sc.Perform(ctx.contract, ctx.client, job.Address, job.ActionMethod, &bind.TransactOpts{
				From:     ctx.signer.From,
				Signer:   ctx.signer.Signer,
				GasLimit: ctx.signer.GasLimit,
			}, nil)

			if err != nil {
				errCounter += 1
				fmt.Println("errCounter: ", errCounter)

				log = fmt.Sprintf("Error while trying to make the tx on performMethod: %v", err)

				if errCounter > maxErrorsLimit {
					fmt.Println("here")
					stopLog := fmt.Sprintf("Failed 3 times so stopped job. Last error: %s", log)
					updateJob(cj.jobstorage, job, stopLog)
					stop <- true
				}

				// Check if the job still doesn't exist, it doesn't need an update
				_, err := cj.jobstorage.GetById(job.ID)
				if err != nil {
					return
				}

				updateJob(cj.jobstorage, job, log)
				return
			}
			// If it succed, the error counter comes back to zero
			errCounter = 0
			fmt.Printf("Tx performed on %s network!: %s \n", job.Network, tx.Hash())
		}
	})

	if err != nil {
		return err
	}

	return nil
}

// TODO(nb): Implement a function that gets and returns the state of the service
func GetState(id string) (state provider.State) {
	return provider.StatusRunning
}

// Method for stopping cronjob when an error is occurred
func updateJob(s *storage.Job, j *job.Job, errorLog string) {
	// Update the status of the job on the storage
	j.Status = provider.StatusError

	// Update the logs field with the error msg
	j.Logs = &errorLog

	fmt.Println("Updating job...")
	s.Update(j)
	fmt.Println("Job updated!")
}
