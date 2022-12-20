package providers

import (
	"github.com/darchlabs/jobs/internal/api"
	"github.com/darchlabs/jobs/internal/provider"
)

type ListProvidersHandler struct {
	providers []provider.Provider
}

func NewListProvidersHandler() *ListProvidersHandler {
	providers := make([]provider.Provider, 0)
	dlNetworks := make([]string, 0)
	chainlinkNetworks := make([]string, 0)

	dlNetworks = append(dlNetworks, "ethereum", "goerli")

	dlKeepers := provider.Provider{
		ID:       "1",
		Name:     "Darch Labs Keepers",
		Networks: dlNetworks,
	}
	chainlinkKeepers := provider.Provider{
		ID:       "2",
		Name:     "Chainlink Keepers",
		Networks: chainlinkNetworks,
	}

	providers = append(providers, dlKeepers, chainlinkKeepers)

	return &ListProvidersHandler{
		providers: providers,
	}
}

func (lp *ListProvidersHandler) Invoke() *api.HandlerRes {
	// Get providers
	providers := lp.providers

	// prepare response
	return &api.HandlerRes{Payload: providers, HttpStatus: 200, Err: nil}
}
