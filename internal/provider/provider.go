package provider

type Provider interface {
	Setup() bool
	Fund(int64) bool

	/// @dev: GetState(Name) (IsWorking, IsSetup, NeedsFunding)
	GetState(string) (bool, bool, bool)
}
