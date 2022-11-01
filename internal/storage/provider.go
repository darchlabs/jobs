package storage

import (
	"encoding/json"

	"github.com/darchlabs/jobs/internal/provider"
)

type Provider struct {
	storage *S
}

func NewProvider(s *S) *Provider {
	return &Provider{
		storage: s,
	}
}

func (p *Provider) List() ([]*provider.Provider, error) {
	data := make([]*provider.Provider, 0)

	// Iterate over the values and append them to the slice
	iter := p.storage.DB.NewIterator(nil, nil)
	for iter.Next() {
		var pp *provider.Provider
		err := json.Unmarshal(iter.Value(), &pp)
		if err != nil {
			return nil, err
		}

		data = append(data, pp)
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return data, nil
}
