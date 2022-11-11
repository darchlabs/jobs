/* Darch Labs implementation for Keepers */

package providermanager

import (
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	jobsabi "github.com/darchlabs/jobs/abi"
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

// TODO(nb): Move to another file
// Method for calling the smart contract view function with bool return
func (m *M) Call(j *job.Job, contract *bind.BoundContract, opts *bind.CallOpts) (*bool, error) {
	var out []interface{}

	abiConn, err := jobsabi.NewAbi(common.HexToAddress(j.Address), m.client)

	fmt.Println(1)
	abiTwo := jobsabi.AbiCallerRaw{Contract: &abiConn.AbiCaller}
	m.checkAndStop(err, j)

	fmt.Println(2)
	err = abiTwo.Call(opts, &out, *j.CheckMethod)
	if err != nil {
		return nil, err
	}

	fmt.Println(3)
	fmt.Println("out: ", out)
	out0 := out[0].(bool)
	fmt.Println("out0: ", out0)

	return &out0, err
}

// TODO(nb): Move to another file independently of Manager
// Method for calling the smart contract view function with bool return
func (m *M) Perform(j *job.Job, contract *bind.BoundContract, opts *bind.TransactOpts, params interface{}) error {
	abiConn, err := jobsabi.NewAbi(common.HexToAddress(j.Address), m.client)
	fmt.Println(1)

	// abiTwo := jobsabi.AbiRaw{Contract: abiConn}

	// abiFive := jobsabi.AbiCallerRaw{Contract: &abiConn.AbiCaller}

	abiThree := jobsabi.AbiTransactorRaw{Contract: &abiConn.AbiTransactor}
	m.checkAndStop(err, j)
	fmt.Println(2)

	fmt.Println(params)
	fmt.Println("params: ", params)
	// fmt.Printf("len(params): %d", params...)

	var a []*interface{}
	a = append(a, nil)

	fmt.Println("a: ", a)
	fmt.Println("len(a): ", len(a))

	if params == nil {
		tx, err := abiThree.Transact(opts, j.ActionMethod)
		if err != nil {
			return err
		}
		fmt.Println(3)
		fmt.Println("tx: ", tx)
		return nil
	}

	tx, err := abiThree.Transact(opts, j.ActionMethod, params)
	fmt.Println("tx: ", tx)
	if err != nil {
		return err
	}

	fmt.Println(3)
	fmt.Println("tx: ", tx)

	return nil
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
		fmt.Println(execute)

		fmt.Println("I'm here!")

		// Get blockchain id
		fmt.Println("Getting network..")
		chainId := getChainId(j.Network)
		if chainId == int64(0) {
			err := fmt.Errorf("invalid chain id for %s network", j.Network)
			m.checkAndStop(err, j)
		}
		fmt.Println("chainId: ", chainId)

		fmt.Println("Getting signer...")
		// Get signer for then execute the tx and evaluate it
		signer, err := getSigner(m.privateKey, *m.client, chainId, nil, nil)
		m.checkAndStop(err, j)
		fmt.Println("signer.From: ", signer.From)

		// Parse address
		address := common.HexToAddress(j.Address)
		fmt.Println("address: ", address)

		// Get an instance of the smart contract
		fmt.Println("client: ", m.client)

		// TODO(nb): Import ca sync abi parser

		abiFormatted := "[{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"perform\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"setStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

		parsedAbi, err := abi.JSON(strings.NewReader(abiFormatted))
		m.checkAndStop(err, j)

		fmt.Println("parsedAbi.Methods", parsedAbi.Methods)
		fmt.Println("Getting contract...")
		contract := GetContract(j.Address, parsedAbi, m.client)
		m.checkAndStop(err, j)

		actionMethod := fmt.Sprintf("%s", parsedAbi.Methods[j.ActionMethod])
		if actionMethod == "" {
			err = fmt.Errorf("there is no %s method inside the contract abi", actionMethod)
			m.checkAndStop(err, j)
		}

		// check if j.CheckMethod is nil. If it is nil execute action methdo directly.
		if j.CheckMethod != nil {
			checkMethod := fmt.Sprintf("%s", parsedAbi.Methods[*j.CheckMethod])
			if checkMethod == "" {
				err = fmt.Errorf("there is no %s method inside the contract abi", checkMethod)
				m.checkAndStop(err, j)
			}
			fmt.Println("checkMethod: ", checkMethod)

			fmt.Println("Checking method...")
			res, err := m.Call(j, contract, &bind.CallOpts{})
			m.checkAndStop(err, j)

			execute = *res
		}

		fmt.Println("execute: ", execute)
		fmt.Println("err: ", err)

		// Execute action method and see return
		if execute && err == nil {
			fmt.Println("1")
			fmt.Println("aa: ", parsedAbi.Methods[j.ActionMethod])

			// Keep printing signer
			fmt.Println("signer.From: ", signer.From)
			fmt.Println("signer.Nonce: ", signer.Nonce)
			fmt.Println("signer.Signer: ", signer.Signer)
			fmt.Println("signer.Value: ", signer.Value)
			fmt.Println("GasPrice:  ", signer.GasPrice)
			fmt.Println("GasFeeCap: ", signer.GasFeeCap)
			fmt.Println("GasTipCap: ", signer.GasTipCap)
			fmt.Println("GasLimit:  ", signer.GasLimit)

			fmt.Println("Performing tx...")
			err = m.Perform(j, contract, &bind.TransactOpts{
				From:      signer.From,
				Nonce:     signer.Nonce,
				Signer:    signer.Signer,
				Value:     signer.Value,
				GasPrice:  signer.GasPrice,
				GasFeeCap: signer.GasFeeCap,
				GasTipCap: signer.GasTipCap,
				GasLimit:  signer.GasLimit,
				Context:   nil,
				NoSend:    false,
			}, nil)
			m.checkAndStop(err, j)

		}
	})

	if err != nil {
		return err
	}

	fmt.Println("Starting cron ...")
	m.cron.Start()
	fmt.Println("Cron started!")

	fmt.Println("Retorning nil...")
	return nil
}

// Method for stopping cronjob when an error is occurred
func (m *M) checkAndStop(err error, j *job.Job) {
	if err != nil {
		fmt.Println("err: ", err)
		fmt.Println("Updating job...")
		j.Status = "error"
		m.jobstorage.Update(j)
		fmt.Println("Job updated!")

		fmt.Println("Stopping cron ...")

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
