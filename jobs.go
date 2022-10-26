package jobs

import prov "github.com/darchlabs/jobs/internal/provider"

type ProviderStorage interface {
	ListImplementations() []*prov.Provider
	GetImplementation(id int8) (*prov.Provider, error)

	/* TODO(nb): Don't know how to add the new imported repo, if post the url to the repo by an API or import in repo and update manually in DB
	 * Maybe, add a hook for when a new implementation is added via API, it is inserted in the DB */
	AddImplementation(implementation *prov.Provider) error
	UpdateImplementation(id int8, implementation *prov.Provider) error
	DeleteImplementation(id int8) error
}

type JobStorage interface {
	GetUsingImplementations(userId string) []*prov.Provider
	GetCurrentContracts(userId string) []string
}
