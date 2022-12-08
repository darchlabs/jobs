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

// ctx for running an already configured cronjob cron
type cronCTX struct {
	signer   *bind.TransactOpts
	abi      abi.ABI
	contract *bind.BoundContract
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

	fmt.Println("Getting signer...")
	signer, err = sc.GetSigner(cj.pk, *cj.client, chainId, nil, nil)
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
	contractCode, err := cj.client.PendingCodeAt(context.Background(), common.HexToAddress(job.Address))
	if err != nil {
		return nil, err
	}

	if len(contractCode) == 0 {
		fmt.Println("yes")
		return nil, fmt.Errorf("%s", "the contract address doesn't exist")
	}

	fmt.Println("Getting contract...")
	contract = sc.GetContract(job.Address, parsedAbi, cj.client)

	// fmt.Println("contract /-->: ", contract)
	if contract == nil {
		return nil, fmt.Errorf("there is no contract under the %s address and abi on the %s network", job.Address, job.Network)
	}
	fmt.Println("Contract ready!")

	// Check the action method exists
	fmt.Println("Getting actionMethod...")
	actionMethod := parsedAbi.Methods[job.ActionMethod].String()
	fmt.Println("actionMethodV0 /-->: ", actionMethod)
	if actionMethod == "" {
		fmt.Println("actionMethod /-->: ", actionMethod)
		return nil, fmt.Errorf("there is no %s method inside the contract abi", actionMethod)
	}
	fmt.Println("actionMethod is OK!")

	fmt.Println("Getting checkMethod...")
	// check if j.CheckMethod is nil. If it is, execute action method directly
	if job.CheckMethod != nil {
		checkMethod := parsedAbi.Methods[*job.CheckMethod].String()
		fmt.Println("checkMethodV0 /-->: ", checkMethod)
		if checkMethod == "" {
			fmt.Println("checkMethod /-->: ", checkMethod)
			return nil, fmt.Errorf("there is no %s method inside the contract abi", checkMethod)
		}
	}
	fmt.Println("checkMethod is OK!")

	return &cronCTX{
		signer:   signer,
		abi:      parsedAbi,
		contract: contract,
	}, err
}

func (cj *Cronjob) SetupAndRun(job *job.Job, ctx *cronCTX) error {
	var err error
	// define variable for the cronjob to know if it must execute or not the job
	execute := true

	// Setup cronjob func
	err = cj.cron.AddFunc(job.Cronjob, func() {
		// if j.CheckMethod is nil, avoid this check
		if job.CheckMethod != nil {
			// Check if the response of the smart contract view function for the cronjob to know if it must execute actionMethod or not
			fmt.Println("Checking method...")
			res, err := sc.Call(ctx.contract, cj.client, job.Address, *job.CheckMethod, &bind.CallOpts{})
			if err != nil {
				err = fmt.Errorf("%s", "Error while trying to call checkMethod")
				cj.updateAndStop(err, job)
				return
			}
			fmt.Println("Check method response: ", *res)

			execute = *res
		}

		// Execute action method and see return if it is needed
		if execute && err == nil {
			fmt.Println("Performing tx...")
			tx, err := sc.Perform(ctx.contract, cj.client, job.Address, job.ActionMethod, &bind.TransactOpts{
				From:     ctx.signer.From,
				Signer:   ctx.signer.Signer,
				GasLimit: ctx.signer.GasLimit,
			}, nil)

			if err != nil {
				err = fmt.Errorf("%s", "Error while trying to make the tx on performMethod")
				cj.updateAndStop(err, job)
				return
			}

			fmt.Println("Tx performed!: ", tx.Hash())
		}
	})

	if err != nil {
		return err
	}

	// Execute cronjob func
	cj.cron.Run()
	fmt.Println("Cron started!")

	j, err := cj.jobstorage.GetById(job.ID)
	if err != nil {
		return err
	}

	if j.Status == string(provider.StatusError) {
		return fmt.Errorf("%s", *j.Logs)
	}

	return nil
}

// TODO(nb): Implement a function that gets and returns the state of the service
func GetState(name string) (state provider.State) {
	return provider.StatusRunning
}

// Method for stopping cronjob when an error is occurred
func (cj *Cronjob) updateAndStop(err error, j *job.Job) {
	fmt.Println("err: ", err)

	// Update the status of the job on the storage
	fmt.Println("Updating job...")
	j.Status = string(provider.StatusError)
	errorLog := err.Error()
	j.Logs = &errorLog
	cj.jobstorage.Update(j)
	fmt.Println("Job updated!")

	// It stops the cronjob and then it'll need to be re-executed to start again
	cj.cron.Stop()
	fmt.Println("Cron stopped!")
}
