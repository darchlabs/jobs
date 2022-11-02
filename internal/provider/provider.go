package provider

// Struct for DB
type Provider struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Networks []string `json:"networks"`
}

// Interface required for each new implementation
type Operator interface {
	Create() error // TODO(nb): define better and implement this
	Setup(c *Config) bool
	GetState(name string) (isWorking bool, isSetup bool, needsFunding bool)
}

// Config struct for setup
type Config struct {
	Address      string
	Abi          string
	CheckMethod  string
	ActionMethod string
	CheckType    string
	CheckValue   string
}
