package jobs

import prov "github.com/darchlabs/jobs/internal/provider"

type ProviderStorage interface {
	List() []*prov.Provider
	Get(id uint8) (*prov.Provider, error)
	Add(implementation *prov.Provider) error
	Update(id uint8, implementation *prov.Provider) error
	DeleteImplementation(id uint8) error
}
