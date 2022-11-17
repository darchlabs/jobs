package storage

import (
	"encoding/json"
	"time"

	"github.com/darchlabs/jobs/internal/job"
	"github.com/teris-io/shortid"
)

type Job struct {
	storage *S
}

func NewJob(s *S) *Job {
	return &Job{
		storage: s,
	}
}

func (j *Job) List() ([]*job.Job, error) {
	data := make([]*job.Job, 0)

	iter := j.storage.DB.NewIterator(nil, nil)
	for iter.Next() {
		var jj *job.Job
		err := json.Unmarshal(iter.Value(), &jj)
		if err != nil {
			return nil, err
		}

		data = append(data, jj)
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (j *Job) Insert(jobInput *job.Job) (*job.Job, error) {
	// generate id for database
	id, err := shortid.Generate()
	if err != nil {
		return nil, err
	}

	jobInput.ID = id
	jobInput.CreatedAt = time.Now()
	// TODO(nb): Create a state struct or pattern available for this field
	jobInput.Status = "idle"

	b, err := json.Marshal(jobInput)
	if err != nil {
		return nil, err
	}

	// save in database
	err = j.storage.DB.Put([]byte(id), b, nil)
	if err != nil {
		return nil, err
	}

	return jobInput, nil
}

func (j *Job) Update(jobInput *job.Job) (*job.Job, error) {
	jobInput.UpdatedAt = time.Now()

	b, err := json.Marshal(jobInput)
	if err != nil {
		return nil, err
	}

	// save in database
	err = j.storage.DB.Put([]byte(jobInput.ID), b, nil)
	if err != nil {
		return nil, err
	}

	return jobInput, nil
}
