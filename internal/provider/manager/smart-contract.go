package providermanager

import (
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
	var out []interface{}

	err := contract.Call(opts, &out, checkMethod)
	if err != nil {
		return nil, err
	}

	out0 := out[0].(bool)
	return &out0, err
}

// Method for calling the smart contract view function with bool return
func Perform(contract *bind.BoundContract, client *ethclient.Client, address string, actionMethod string, opts *bind.TransactOpts, params interface{}) (*types.Transaction, error) {
	if params == nil {
		tx, err := contract.Transact(opts, actionMethod)
		if err != nil {
			return nil, err
		}

		return tx, nil
	}

	tx, err := contract.Transact(opts, actionMethod, params)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
