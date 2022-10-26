package providerstorage

import (
	"encoding/json"
	"fmt"

	"github.com/darchlabs/jobs/internal/provider"
	"github.com/darchlabs/jobs/internal/storage"
)

type ProviderStorage struct {
	storage *storage.S
}

func New(s *storage.S) *ProviderStorage {
	return &ProviderStorage{
		storage: s,
	}
}

func (ps *ProviderStorage) AddImplementation(p *provider.Provider) error {
	// Validate if the provider already exists
	iter := ps.storage.DB.NewIterator(nil, nil)
	for iter.Next() {
		var pprovider *provider.Provider
		err := json.Unmarshal(iter.Value(), &pprovider)

		if err != nil {
			return err
		}

		if p.Name == pprovider.Name {
			return fmt.Errorf("%s", "This provider already exists")
		}

		implType := fmt.Sprintf("%T", p.Implemetation)
		if implType != "provider.Implementation" {
			return fmt.Errorf("%s", "The implementation doesn't fit the Provider interface required")
		}

	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		return err
	}

	// TODO(nb): Generate and set id incremental
	p.Id = "1"

	// JSON.stringify
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	// Save in database
	ps.storage.DB.Put([]byte(p.Id), b, nil)
	if err != nil {
		return err
	}

	return nil
}
