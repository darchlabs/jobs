package provider

// Struct for DB
type Provider struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Networks []string `json:"networks,omitempty"`
}

// Interface required for each new implementation
type Implementation interface {
	Setup(address string, abi string, checkMethod string, actionMethod string, checkType string, checkValue string) bool
	GetState(name string) (IsWorking bool, IsSetup bool, NeedsFunding bool)
	// Fund(amount uint64) bool
}
