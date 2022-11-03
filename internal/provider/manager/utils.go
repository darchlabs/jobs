package providermanager

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TODO(nb): ask ca if the nonce is necessary or not to obtain and why?
func getSigner(pk string, client ethclient.Client, chainId int64, gasPrice *int64, gasLimit *uint64) (*bind.TransactOpts, error) {
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

func getChainId(name string) int64 {
	// TODO(nb): hardcode all the chain id for the chains that'll be used
	networksMap := map[string]int64{
		"ethereum": int64(1),
	}

	return networksMap[name]
}
