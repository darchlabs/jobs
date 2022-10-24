package provider

type Provider interface {
	Setup() bool
	Fund(amount int64) bool
	GetState(name string) (IsWorking bool, IsSetup bool, NeedsFunding bool)
}
