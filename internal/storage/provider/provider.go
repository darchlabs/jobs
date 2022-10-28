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

/* Add, Update and Delete */
// func (ps *PS) AddImplementation(p *provider.Provider) error {
// 	var id uint8
// 	// Validate the implementation received is the correct one
// 	// implType := fmt.Sprintf("%T", provider.Implementation)
// 	// if implType != "provider.Implementation" {
// 	// 	return fmt.Errorf("%s", "The implementation doesn't fit the Provider interface required")
// 	// }

// 	// Validate if the provider already exists
// 	iter := ps.storage.DB.NewIterator(nil, nil)
// 	for iter.Next() {
// 		var pprovider *provider.Provider
// 		err := json.Unmarshal(iter.Value(), &pprovider)

// 		if err != nil {
// 			return err
// 		}

// 		if p.Name == pprovider.Name {
// 			return fmt.Errorf("%s", "This provider already exists")
// 		}

// 	}
// 	iter.Release()

// 	err := iter.Error()
// 	if err != nil {
// 		return err
// 	}

// 	// add 1 to the last id in order to make it incremental
// 	p.Id = id

// 	// JSON.stringify
// 	b, err := json.Marshal(p)
// 	if err != nil {
// 		return err
// 	}

// 	// Save in database
// 	err = ps.storage.DB.Put([]byte(fmt.Sprintf("%d", p.Id)), b, nil)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
// func (ps *PS) UpdateImplementation(id uint8, p *provider.Provider) error {
// 	// Validate the implementation exists
// 	data, err := ps.storage.DB.Get([]byte(fmt.Sprintf("%d", id)), nil)
// 	if err != nil {
// 		return err
// 	}

// 	if data == nil {
// 		return fmt.Errorf("%s", "No provider exists for the given id")
// 	}

// 	// Validate the provider implementation received is crrect
// 	implType := fmt.Sprintf("%T", p.Implemetation)
// 	if implType != "provider.Implementation" {
// 		return fmt.Errorf("%s", "The implementation doesn't fit the Provider interface required")
// 	}

// 	// Convert provider to bytes
// 	b, err := json.Marshal(p)
// 	if err != nil {
// 		return err
// 	}

// 	// Update in DB
// 	err = ps.storage.DB.Put([]byte(fmt.Sprintf("%d", id)), b, nil)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (ps *PS) DeleteImplementation(id uint8) error {
// 	// Delete provider from DB
// 	err := ps.storage.DB.Delete([]byte(fmt.Sprintf("%d", id)), nil)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
