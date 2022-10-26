package providerstorage

import (
	"encoding/json"
	"fmt"
	"log"

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

// TODO(nb): add authorization methods only for Darch Labs team IN HANDLERS

func (ps *PS) AddImplementation(p *provider.Provider) error {
	// Validate the implementation received is the correct one
	implType := fmt.Sprintf("%T", p.Implemetation)
	if implType != "provider.Implementation" {
		return fmt.Errorf("%s", "The implementation doesn't fit the Provider interface required")
	}

	// Validate if the provider already exists and get last id
	var id int8
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

		// TODO(nb): iteration over leveldb respects the order? The last item iterated is indeed the last inserted?
		// Get the last id
		id = pprovider.Id

	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		return err
	}

	// add 1 to the last id in order to make it incremental
	p.Id = id + 1

	// JSON.stringify
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	// Save in database
	err = ps.storage.DB.Put([]byte(fmt.Sprintf("%d", p.Id)), b, nil)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PS) ListImplementations() []*provider.Provider {
	providers := make([]*provider.Provider, 0)

	iter := ps.storage.DB.NewIterator(nil, nil)
	for iter.Next() {
		var p *provider.Provider
		err := json.Unmarshal(iter.Value(), &p)

		if err != nil {
			log.Printf("Error while iterating db: %v \n", err)
			return nil
		}

		providers = append(providers, p)
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		log.Printf("Error while iterating db: %v \n", err)
		return nil
	}

	return providers
}

func (ps *PS) GetImplementation(id int8) (*provider.Provider, error) {
	// Get the provider by the id key in bytes
	data, err := ps.storage.DB.Get([]byte(fmt.Sprintf("%d", id)), nil)
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

func (ps *PS) UpdateImplementation(id int8, p *provider.Provider) error {
	// Validate the implementation exists
	data, err := ps.storage.DB.Get([]byte(fmt.Sprintf("%d", id)), nil)
	if err != nil {
		return err
	}

	if data == nil {
		return fmt.Errorf("%s", "No provider exists for the given id")
	}

	// Validate the provider implementation received is crrect
	implType := fmt.Sprintf("%T", p.Implemetation)
	if implType != "provider.Implementation" {
		return fmt.Errorf("%s", "The implementation doesn't fit the Provider interface required")
	}

	// Convert provider to bytes
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	// Update in DB
	err = ps.storage.DB.Put([]byte(fmt.Sprintf("%d", id)), b, nil)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PS) DeleteImplementation(id int8) error {
	// Delete provider from DB
	err := ps.storage.DB.Delete([]byte(fmt.Sprintf("%d", id)), nil)
	if err != nil {
		return err
	}

	return nil
}
