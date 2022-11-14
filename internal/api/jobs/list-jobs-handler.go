package jobsapi

import (
	"github.com/darchlabs/jobs/internal/api"
	"github.com/darchlabs/jobs/internal/storage"
)

type ListJobsHandler struct {
	storage *storage.Job
}

func NewListJobsHandler(js *storage.Job) *ListJobsHandler {
	return &ListJobsHandler{
		storage: js,
	}
}

func (ListJobsHandler) Invoke(ctx Context) *api.HandlerRes {
	// Get elements from db
	data, err := ctx.JobStorage.List()
	if err != nil {
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}

	// Prepare response
	return &api.HandlerRes{Payload: data, HttpStatus: 200, Err: nil}
}
