package providermanager

import "github.com/darchlabs/jobs/internal/job"

// TODO(nb): V2 create synchronizer based keeper code for exec a smart ocntract by an event listening
type Synchronizer struct {
}

func NewSynchronizer() *Synchronizer {
	return &Synchronizer{}
}

// Implementation when a listener over events is needed (synchronizer) --> Jobs V2
func (s *Synchronizer) SetupAndRun(j *job.Job) error {
	return nil
}
