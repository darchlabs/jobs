package provider

// Struct for DB
type Provider struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Networks []string `json:"networks"`
}

// Interface required for each new implementation
type Operator interface {
<<<<<<< HEAD
	Setup(c *Config) error
	GetState(name string) (state State)
=======
	Create() error // TODO(nb): define better and implement this
	Setup(c *Config) bool
	GetState(name string) (isWorking bool, isSetup bool, needsFunding bool)
>>>>>>> nb-feat/create-jobs-endpoint
}

// Config struct for setup
type Config struct {
	Provider     Provider
	Address      string
	Abi          string
	CheckMethod  string
	ActionMethod string
	CheckType    string
	CheckValue   string
}

type State string

const (
	StatusIdle     State = "idle"
	StatusStarting State = "starting"
	StatusRunning  State = "running"
	StatusStopping State = "stopping"
	StatusStopped  State = "stopped"
	StatusError    State = "error"
)
