package providers

import providerstorage "github.com/darchlabs/jobs/internal/storage/provider"

type ListProvidersHandler struct {
	storage providerstorage.PS
}

func NewListProvidersHandler(ps providerstorage.PS) *ListProvidersHandler {
	return &ListProvidersHandler{
		storage: ps,
	}
}

func (ListProvidersHandler) Invoke(ctx Context) *handlerRes {
	// Get elements from db
	data, err := ctx.ProviderStorage.List()
	if err != nil {
		return &handlerRes{err.Error(), 500, err}
	}

	return &handlerRes{data, 200, nil}
}
