package providermanager

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// // TokenMetaData contains all meta data concerning the Token contract.
// var TokenMetaData = &bind.MetaData{
// 	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"Balance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amt\",\"type\":\"uint256\"}],\"name\":\"Deposite\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amt\",\"type\":\"uint256\"}],\"name\":\"Withdrawl\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
// 	Bin: "0x60806040526000805534801561001457600080fd5b50600180546001600160a01b0319163317905561002f610034565b610072565b34600080828254610045919061004c565b9091555050565b6000821982111561006d57634e487b7160e01b600052601160045260246000fd5b500190565b6101e5806100816000396000f3fe6080604052600436106100435760003560e01c80630ef678871461005757806342002bc81461007b578063e615c1a01461009b578063f851a440146100bb57600080fd5b36610052576100506100f3565b005b600080fd5b34801561006357600080fd5b506000545b6040519081526020015b60405180910390f35b34801561008757600080fd5b50610068610096366004610151565b61010b565b3480156100a757600080fd5b506100506100b6366004610151565b610126565b3480156100c757600080fd5b506001546100db906001600160a01b031681565b6040516001600160a01b039091168152602001610072565b346000808282546101049190610180565b9091555050565b60008160005461011b9190610180565b600081905592915050565b6001546001600160a01b0316331461013d57600080fd5b8060005461014b9190610198565b60005550565b60006020828403121561016357600080fd5b5035919050565b634e487b7160e01b600052601160045260246000fd5b600082198211156101935761019361016a565b500190565b6000828210156101aa576101aa61016a565b50039056fea26469706673582212202a886a16db3a2519c47f56f6c0a9aa72d1ffa8b18909f2c5e7edebc433eceaab64736f6c634300080a0033",
// }

// // TokenABI is the input ABI used to generate the binding from.
// // Deprecated: Use TokenMetaData.ABI instead.
// var TokenABI = TokenMetaData.ABI

// // TokenBin is the compiled bytecode used for deploying new contracts.
// // Deprecated: Use TokenMetaData.Bin instead.
// var TokenBin = TokenMetaData.Bin

// DeployToken deploys a new Ethereum contract, binding an instance of Token to it.
func DeployToken(auth *bind.TransactOpts, tokenMetaData bind.MetaData, backend bind.ContractBackend) (common.Address, *types.Transaction, *Token, error) {
	parsed, err := tokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(tokenMetaData.Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Token{TokenCaller: TokenCaller{contractCaller: contract}, TokenTransactor: TokenTransactor{contractTransactor: contract}, TokenFilterer: TokenFilterer{ContractFilterer: contract}}, nil
}

// Token is an auto generated Go binding around an Ethereum contract.
type Token struct {
	TokenCaller     // Read-only binding to the contract
	TokenTransactor // Write-only binding to the contract
	TokenFilterer   // Log filterer for contract events
}

// TokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCaller struct {
	contractCaller *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenTransactor struct {
	contractTransactor *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenFilterer struct {
	ContractFilterer *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenSession struct {
	Contract     *Token            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCallerSession struct {
	Contract *TokenCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenTransactorSession struct {
	Contract     *TokenTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenRaw struct {
	Contract *Token // Generic contract binding to access the raw methods on
}

// TokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCallerRaw struct {
	Contract *TokenCaller // Generic read-only contract binding to access the raw methods on
}

// TokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenTransactorRaw struct {
	Contract *TokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewToken creates a new instance of Token, bound to a specific deployed contract.
func NewToken(address common.Address, abi string, backend bind.ContractBackend) (*Token, error) {
	contract, err := bindToken(address, abi, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Token{TokenCaller: TokenCaller{contractCaller: contract}, TokenTransactor: TokenTransactor{contractTransactor: contract}, TokenFilterer: TokenFilterer{ContractFilterer: contract}}, nil
}

// NewTokenCaller creates a new read-only instance of Token, bound to a specific deployed contract.
func NewTokenCaller(address common.Address, abi string, caller bind.ContractCaller) (*TokenCaller, error) {
	contract, err := bindToken(address, abi, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCaller{contractCaller: contract}, nil
}

// NewTokenTransactor creates a new write-only instance of Token, bound to a specific deployed contract.
func NewTokenTransactor(address common.Address, abi string, transactor bind.ContractTransactor) (*TokenTransactor, error) {
	contract, err := bindToken(address, abi, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenTransactor{contractTransactor: contract}, nil
}

// NewTokenFilterer creates a new log filterer instance of Token, bound to a specific deployed contract.
func NewTokenFilterer(address common.Address, abi string, filterer bind.ContractFilterer) (*TokenFilterer, error) {
	contract, err := bindToken(address, abi, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenFilterer{ContractFilterer: contract}, nil
}

// bindToken binds a generic wrapper to an already deployed contract.
func bindToken(address common.Address, tokenAbi string, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(tokenAbi))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token *TokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Token.Contract.TokenCaller.contractCaller.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token *TokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.Contract.TokenTransactor.contractTransactor.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token *TokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token.Contract.TokenTransactor.contractTransactor.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token *TokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Token.Contract.contractCaller.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token *TokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.Contract.contractTransactor.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token *TokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token.Contract.contractTransactor.Transact(opts, method, params...)
}

// Balance is a free data retrieval call binding the contract method 0x0ef67887.
//
// Solidity: function Balance() view returns(uint256)
func (_Token *TokenCaller) Balance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Token.contractCaller.Call(opts, &out, "Balance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Balance is a free data retrieval call binding the contract method 0x0ef67887.
//
// Solidity: function Balance() view returns(uint256)
func (_Token *TokenSession) Balance() (*big.Int, error) {
	return _Token.Contract.Balance(&_Token.CallOpts)
}

// Balance is a free data retrieval call binding the contract method 0x0ef67887.
//
// Solidity: function Balance() view returns(uint256)
func (_Token *TokenCallerSession) Balance() (*big.Int, error) {
	return _Token.Contract.Balance(&_Token.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Token *TokenCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Token.contractCaller.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Token *TokenSession) Admin() (common.Address, error) {
	return _Token.Contract.Admin(&_Token.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Token *TokenCallerSession) Admin() (common.Address, error) {
	return _Token.Contract.Admin(&_Token.CallOpts)
}

// Deposite is a paid mutator transaction binding the contract method 0x42002bc8.
//
// Solidity: function Deposite(uint256 amt) returns(uint256)
func (_Token *TokenTransactor) Deposite(opts *bind.TransactOpts, amt *big.Int) (*types.Transaction, error) {
	return _Token.contractTransactor.Transact(opts, "Deposite", amt)
}

// Deposite is a paid mutator transaction binding the contract method 0x42002bc8.
//
// Solidity: function Deposite(uint256 amt) returns(uint256)
func (_Token *TokenSession) Deposite(amt *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Deposite(&_Token.TransactOpts, amt)
}

// Deposite is a paid mutator transaction binding the contract method 0x42002bc8.
//
// Solidity: function Deposite(uint256 amt) returns(uint256)
func (_Token *TokenTransactorSession) Deposite(amt *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Deposite(&_Token.TransactOpts, amt)
}

// Withdrawl is a paid mutator transaction binding the contract method 0xe615c1a0.
//
// Solidity: function Withdrawl(uint256 _amt) returns()
func (_Token *TokenTransactor) Withdrawl(opts *bind.TransactOpts, _amt *big.Int) (*types.Transaction, error) {
	return _Token.contractTransactor.Transact(opts, "Withdrawl", _amt)
}

// Withdrawl is a paid mutator transaction binding the contract method 0xe615c1a0.
//
// Solidity: function Withdrawl(uint256 _amt) returns()
func (_Token *TokenSession) Withdrawl(_amt *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Withdrawl(&_Token.TransactOpts, _amt)
}

// Withdrawl is a paid mutator transaction binding the contract method 0xe615c1a0.
//
// Solidity: function Withdrawl(uint256 _amt) returns()
func (_Token *TokenTransactorSession) Withdrawl(_amt *big.Int) (*types.Transaction, error) {
	return _Token.Contract.Withdrawl(&_Token.TransactOpts, _amt)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Token *TokenTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.contractTransactor.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Token *TokenSession) Receive() (*types.Transaction, error) {
	return _Token.Contract.Receive(&_Token.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Token *TokenTransactorSession) Receive() (*types.Transaction, error) {
	return _Token.Contract.Receive(&_Token.TransactOpts)
}
