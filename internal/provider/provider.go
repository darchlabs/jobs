package provider

// Struct for DB
type Provider struct {
	Id            string
	Name          string
	Implemetation Implementation
}

// Interface required for each new implementation
type Implementation interface {
	Setup() bool
	Fund(amount int64) bool
	GetState(name string) (IsWorking bool, IsSetup bool, NeedsFunding bool)
}
