package providerstorage

import "github.com/darchlabs/jobs/internal/provider"

type ProviderStorage struct {
	Id            string
	Name          string
	Implemetation provider.Provider
}
