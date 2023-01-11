package smartcontracts

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

func GetSigner(pk string, client ethclient.Client, chainId int64, gasPrice *int64, gasLimit *uint64) (*bind.TransactOpts, error) {
	// parse private key to ECDSA
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, err
	}

	// set gas price param
	var gp *big.Int
	if gasPrice != nil {
		gp = big.NewInt(*gasPrice)
	} else {
		// get recommended gas price
		gp, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, err
		}
	}

	// set gas limit param (in units)
	var gl uint64
	if gasLimit != nil {
		gl = *gasLimit
	} else {
		gl = uint64(300000)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(chainId)))
	if err != nil {
		return nil, err
	}

	auth.GasPrice = gp
	auth.GasLimit = gl

	return auth, nil
}
