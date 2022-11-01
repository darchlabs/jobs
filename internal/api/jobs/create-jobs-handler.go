package jobsapi

import (
	"encoding/json"

	"github.com/darchlabs/jobs/internal/api"
	"github.com/darchlabs/jobs/internal/job"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/go-playground/validator"
)

type CreateJobsHandler struct {
	storage storage.Job
}

func NewCreateJobsHandler(js storage.Job) *CreateJobsHandler {
	return &CreateJobsHandler{
		storage: js,
	}
}

func (CreateJobsHandler) Invoke(ctx Context) *api.HandlerRes {
	// Prepare body request struct
	body := struct {
		Job *job.Job `json:"job"`
	}{}

	// Parse body to Job struct
	err := json.Unmarshal(ctx.c.Body(), &body)
	if err != nil {
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}

	validate := validator.New()
	validate.Struct(ctx.JobStorage)

	// Insert job in jobstorage DB
	j, err := ctx.JobStorage.Insert(body.Job)
	if err != nil {
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}

	return &api.HandlerRes{Payload: j, HttpStatus: 200, Err: nil}
}
