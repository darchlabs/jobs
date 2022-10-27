package provider

// Struct for DB
type Provider struct {
	Id            uint8
	Name          string
	Implemetation Implementation
}

// Interface required for each new implementation
type Implementation interface {
	Setup(scAddress string) bool
	Fund(amount uint64) bool
	GetState(name string) (IsWorking bool, IsSetup bool, NeedsFunding bool)
}
