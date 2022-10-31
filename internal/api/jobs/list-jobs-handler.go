package jobsapi

import (
	jobstorage "github.com/darchlabs/jobs/internal/storage/job"
)

type ListJobsHandler struct {
	storage jobstorage.JS
}

func NewListJobsHandler(js jobstorage.JS) *ListJobsHandler {
	return &ListJobsHandler{
		storage: js,
	}
}

func (ListJobsHandler) Invoke(ctx Context) *HandlerRes {
	// Get elements from db
	data, err := ctx.JobStorage.List()
	if err != nil {
		return &HandlerRes{err.Error(), 500, err}
	}

	// Prepare response
	return &HandlerRes{data, 200, nil}
}
