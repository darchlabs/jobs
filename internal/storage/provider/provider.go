package providerstorage

import (
	"encoding/json"

	"github.com/darchlabs/jobs/internal/provider"
	"github.com/darchlabs/jobs/internal/storage"
)

type PS struct {
	storage *storage.S
}

func New(s *storage.S) *PS {
	return &PS{
		storage: s,
	}
}

func (ps *PS) List() ([]*provider.Provider, error) {
	data := make([]*provider.Provider, 0)

	// Iterate over the values and append them to the slice
	iter := ps.storage.DB.NewIterator(nil, nil)
	for iter.Next() {
		var p *provider.Provider
		err := json.Unmarshal(iter.Value(), &p)
		if err != nil {
			return nil, err
		}

		data = append(data, p)
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ps *PS) GetImplementation(id string) (*provider.Provider, error) {
	// Get the provider by the id key in bytes
	data, err := ps.storage.DB.Get([]byte(id), nil)
	if err != nil {
		return nil, err
	}

	// Parse the bytes to provider struct type
	var p *provider.Provider

	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
