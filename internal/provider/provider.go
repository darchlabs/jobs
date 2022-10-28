package provider

// Struct for DB
type Provider struct {
	Id       string
	Name     string
	Networks []string
}

// Interface required for each new implementation
type Implementation interface {
	Setup(address string, abi string, checkMethod string, actionMethod string, checkType string, checkValue string) bool
	GetState(name string) (IsWorking bool, IsSetup bool, NeedsFunding bool)
	// Fund(amount uint64) bool
}
