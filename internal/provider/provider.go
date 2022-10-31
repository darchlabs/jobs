package provider

// Struct for DB
type Provider struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Networks []string `json:"networks"`
}

// Interface required for each new implementation
type Create interface {
	Setup(address string, abi string, checkMethod string, actionMethod string, checkType string, checkValue string) bool
	GetState(name string) (isWorking bool, isSetup bool, needsFunding bool)
}
