package jobs

import prov "github.com/darchlabs/jobs/internal/provider"

type ProviderStorage interface {
	ListImplementations() []*prov.Provider
	GetImplementation(id uint8) (*prov.Provider, error)
	AddImplementation(implementation *prov.Provider) error
	UpdateImplementation(id uint8, implementation *prov.Provider) error
	DeleteImplementation(id uint8) error
}

// methods for interacting w DB that intersects User and Provider tables
type JobStorage interface {
	GetUsingImplementations(userId string) []*prov.Provider
	GetCurrentContracts(userId string) []string
	AddUserProvider(scAddress string, event string, userId string, providerId string) error
	UpdateUserProvider(scAddress string, providerId string, needsFunding bool, working bool) error
}
