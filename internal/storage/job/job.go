package jobstorage

import (
	"encoding/json"
	"time"

	"github.com/darchlabs/jobs/internal/job"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/teris-io/shortid"
)

type JS struct {
	storage *storage.S
}

func New(s *storage.S) *JS {
	return &JS{
		storage: s,
	}
}

func (js *JS) List() ([]*job.Job, error) {
	data := make([]*job.Job, 0)

	iter := js.storage.DB.NewIterator(nil, nil)
	for iter.Next() {
		var j *job.Job
		err := json.Unmarshal(iter.Value(), &j)
		if err != nil {
			return nil, err
		}

		data = append(data, j)
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (js *JS) Insert(j *job.Job) (*job.Job, error) {
	// generate id for database
	id, err := shortid.Generate()
	if err != nil {
		return nil, err
	}

	j.ID = id
	j.CreatedAt = time.Now()
	// TODO(nb): Create a state struct or pattern available for this field
	j.Status = "idle"

	b, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}

	// save in database
	err = js.storage.DB.Put([]byte(id), b, nil)
	if err != nil {
		return nil, err
	}

	return j, nil
}
