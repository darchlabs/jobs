package provider

type Provider struct {
	Name         string
	NeedsFunding bool
	Active       bool
	IsSetup      bool
	IsWorking    bool

	Setup func() bool
	Fund  func(int64) bool

	/// @dev: GetState(Name) (IsWorking, IsSetup, NeedsFunding)
	GetState func(string) (bool, bool, bool)
}
