package provider

// Struct for DB
type Provider struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Networks []string `json:"networks"`
}

// Interface required for each new implementation
type Operator interface {
	Setup(c *Config) error
	GetState(name string) (state State)
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
