package providermanager

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func GetContract(address string, inputAbi abi.ABI, backend bind.ContractBackend) *bind.BoundContract {
	contract := bind.NewBoundContract(common.HexToAddress(address), inputAbi, backend, backend, backend)
	return contract
}
