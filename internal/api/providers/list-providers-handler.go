package providers

import (
	"github.com/darchlabs/jobs/internal/api"
	"github.com/darchlabs/jobs/internal/storage"
)

type ListProvidersHandler struct {
	storage storage.Provider
}

func NewListProvidersHandler(ps storage.Provider) *ListProvidersHandler {
	return &ListProvidersHandler{
		storage: ps,
	}
}

func (ListProvidersHandler) Invoke(ctx Context) *api.HandlerRes {
	// Get elements from db
	data, err := ctx.ProviderStorage.List()
	if err != nil {
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}

	// prepare response
	return &api.HandlerRes{Payload: data, HttpStatus: 200, Err: nil}
}
