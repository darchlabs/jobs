package jobsapi

import (
	"encoding/json"
	"fmt"

	"github.com/darchlabs/jobs/internal/api"
	"github.com/darchlabs/jobs/internal/job"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/go-playground/validator"
	"github.com/robfig/cron"
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

	if body.Job.Type == "cronjob" {
		// Validate the cronjob received is correct
		specParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		_, err = specParser.Parse(body.Job.Cronjob)

		if err != nil {
			fmt.Println(err)
			return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
		}

		// Execute manager in order to execute the job
		err = ctx.Manager.Create(body.Job)
		if err != nil {
			return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
		}
	}

	// Insert job in jobstorage DB
	j, err := ctx.JobStorage.Insert(body.Job)
	if err != nil {
		fmt.Println(err)
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}

	return &api.HandlerRes{Payload: j, HttpStatus: 200, Err: nil}
}
