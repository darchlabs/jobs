package providermanager

import (
	jobsabi "github.com/darchlabs/jobs/abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetContract(address string, inputAbi abi.ABI, backend bind.ContractBackend) *bind.BoundContract {
	contract := bind.NewBoundContract(common.HexToAddress(address), inputAbi, backend, backend, backend)
	return contract
}

// Method for calling the smart contract view function with bool return
func Call(contract *bind.BoundContract, client *ethclient.Client, address string, checkMethod string, opts *bind.CallOpts) (*bool, error) {
	abiConn, err := jobsabi.NewAbiCaller(common.HexToAddress(address), client)
	if err != nil {
		return nil, err
	}

	abiCaller := jobsabi.AbiCallerRaw{Contract: abiConn}

	var out []interface{}
	err = abiCaller.Call(opts, &out, checkMethod)
	if err != nil {
		return nil, err
	}

	out0 := out[0].(bool)
	return &out0, err
}

// Method for calling the smart contract view function with bool return
func Perform(contract *bind.BoundContract, client *ethclient.Client, address string, actionMethod string, opts *bind.TransactOpts, params interface{}) (*types.Transaction, error) {
	abiConn, err := jobsabi.NewAbiTransactor(common.HexToAddress(address), client)
	if err != nil {
		return nil, err
	}

	abiTransactor := jobsabi.AbiTransactorRaw{Contract: abiConn}

	if params == nil {
		tx, err := abiTransactor.Transact(opts, actionMethod)
		if err != nil {
			return nil, err
		}

		return tx, nil
	}

	tx, err := abiTransactor.Transact(opts, actionMethod, params)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
