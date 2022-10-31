package jobsapi

import (
	"encoding/json"

	"github.com/darchlabs/jobs/internal/job"
	jobstorage "github.com/darchlabs/jobs/internal/storage/job"
	"github.com/go-playground/validator"
)

type CreateJobsHandler struct {
	storage jobstorage.JS
}

func NewCreateJobsHandler(js jobstorage.JS) *CreateJobsHandler {
	return &CreateJobsHandler{
		storage: js,
	}
}

func (CreateJobsHandler) Invoke(ctx Context) *HandlerRes {
	// Prepare body request struct
	body := struct {
		Job *job.Job `json:"job"`
	}{}

	// Parse body to Job struct
	err := json.Unmarshal(ctx.c.Body(), &body)
	if err != nil {
		return &HandlerRes{err.Error(), 500, err}
	}

	// TODO(nb): validate
	validate := validator.New()
	validate.Struct(ctx.JobStorage)

	// Insert job in jobstorage DB
	j, err := ctx.JobStorage.Insert(body.Job)
	if err != nil {
		return &HandlerRes{err.Error(), 500, err}
	}

	return &HandlerRes{j, 200, nil}
}
