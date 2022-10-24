package jobstorage

// Main DB table

// TODO(nb): The fields should be an id or the struct?
// TODO(nb): How to link it with storage?

type Job struct {
	UserId          string
	ProviderId      string
	SmartContractId string
	SynchronizerId  string
}
