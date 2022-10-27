package implementation

import (
	"github.com/darchlabs/jobs/internal/provider"
	providerstorage "github.com/darchlabs/jobs/internal/storage/provider"
)

// Interface w methods of rinteracting w providers implementations and the ProviderStorage DB
type ProviderImplementation interface {
	AddImplementation(ps *providerstorage.PS, p *provider.Provider) error
	UpdateImplementation(ps *providerstorage.PS, id uint8, p *provider.Provider) error
	DeleteImplementation(ps *providerstorage.PS, id uint8) error
}
