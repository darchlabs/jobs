package jobs

import prov "github.com/darchlabs/jobs/internal/provider"

type ImplementationStorage interface {
	ListImplementations() []*prov.Provider
	GetImplementation(implementationName string) (*prov.Provider, error)
	AddImplementation(implementation *prov.Provider) error
	UpdateImplementation(implementation *prov.Provider) error
	DeleteImplementation(implementation *prov.Provider) error
}

type UserJobsStorage interface {
	GetUsingImplementations(userId string) []*prov.Provider
	GetCurrentContracts(userId string) []string
}
