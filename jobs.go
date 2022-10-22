package jobs

import impl "github.com/darchlabs/jobs/internal/implementations"

type ImplementationStorage interface {
	ListImplementations() []*impl.Implementation
	GetImplementation(implementationName string) (*impl.Implementation, error)
	AddImplementation(implementation *impl.Implementation) error
	UpdateImplementation(implementation *impl.Implementation) error
	DeleteImplementation(implementation *impl.Implementation) error
}

type UserJobsStorage interface {
	GetUsingImplementations(userId string) []*impl.Implementation
	GetCurrentContracts(userId string) []string
}
