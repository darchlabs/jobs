package jobstorage

import "github.com/darchlabs/jobs/internal/storage"

// Main DB table

type Storage struct {
	storage *storage.S
}

func New(s *storage.S) *Storage {
	return &Storage{
		storage: s,
	}
}

// TODO(nb): The fields should be an id or the struct?
type JobStorage struct {
	SmartContractId string
	Event           string
	UserId          string
	SynchronizerId  string
	ProviderId      string

	/// @notice: Provider-User related fields
	Selected     bool
	Setup        bool
	Working      bool
	NeedsFunding bool
}
