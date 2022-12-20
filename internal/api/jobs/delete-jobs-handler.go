package jobsapi

import (
	"fmt"

	"github.com/darchlabs/jobs/internal/api"
	"github.com/darchlabs/jobs/internal/provider"
	"github.com/darchlabs/jobs/internal/storage"
)

type DeleteJobHandler struct {
	storage *storage.Job
}

func NewDeleteJobHandler(js *storage.Job) *DeleteJobHandler {
	return &DeleteJobHandler{
		storage: js,
	}
}

func (DeleteJobHandler) Invoke(ctx Context) *api.HandlerRes {
	// Get id param and assert is not empty
	id := ctx.c.Params("id")
	if id == "" {
		err := fmt.Errorf("%s", "id param in route is empty")
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}

	// Get and check that the job exists in DB
	job, err := ctx.JobStorage.GetById(id)
	if err != nil {
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}

	// Stop job if is running
	if job.Status == provider.StatusRunning {
		ctx.Manager.Stop(id)
	}

	// Delete job from jobstorage DB
	err = ctx.JobStorage.Delete(id)
	if err != nil {
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}

	res := map[string]interface{}{"id": job.ID, "status": "Deleted"}
	return &api.HandlerRes{Payload: res, HttpStatus: 200, Err: nil}
}
