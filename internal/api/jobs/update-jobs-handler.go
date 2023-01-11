package jobsapi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/darchlabs/jobs/internal/api"
	"github.com/darchlabs/jobs/internal/provider"
	"github.com/darchlabs/jobs/internal/storage"
	"github.com/go-playground/validator"
)

type UpdateBody struct {
	Network    string `json:"network"`
	NodeURL    string `json:"nodeUrl"`
	Privatekey string `json:"privateKey"`
	Address    string `json:"address"`
	Cronjob    string `json:"cronjob"`
}

type UpdateJobHandler struct {
	storage *storage.Job
}

func NewUpdateJobHandler(js *storage.Job) *UpdateJobHandler {
	return &UpdateJobHandler{
		storage: js,
	}
}

func (UpdateJobHandler) Invoke(ctx Context) *api.HandlerRes {
	id := ctx.c.Params("id")
	if id == "" {
		err := fmt.Errorf("%s", "id param in route is empty")
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}

	body := struct {
		Req *UpdateBody `json:"job"`
	}{}

	// Parse body to Job struct
	err := json.Unmarshal(ctx.c.Body(), &body)
	if err != nil {
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}

	fmt.Println("Validating...")
	validate := validator.New()
	validate.Struct(ctx.JobStorage)
	fmt.Println("Validated!")

	fmt.Println("Getting job...")
	job, err := ctx.JobStorage.GetById(id)
	if err != nil {
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}
	fmt.Println("Job got!")

	// This is for checking that at least there is a param to be modified in the request
	empty := true

	// Update jobs params that were sent in the request
	if body.Req.Address != "" && body.Req.Address != job.Address {
		job.Address = body.Req.Address
		empty = false
	}

	if body.Req.Cronjob != "" && body.Req.Cronjob != job.Cronjob {
		job.Cronjob = body.Req.Cronjob
		empty = false
	}

	if body.Req.Network != "" && body.Req.Network != job.Network {
		job.Network = body.Req.Network
		empty = false
	}

	if body.Req.NodeURL != "" && body.Req.NodeURL != job.NodeURL {
		job.NodeURL = body.Req.NodeURL
		empty = false
	}

	if body.Req.Privatekey != "" && body.Req.Privatekey != job.Privatekey {
		job.Privatekey = body.Req.Privatekey
		empty = false
	}

	// If empty is true, it means that no params will be updated so an error is returned
	if empty {
		err = fmt.Errorf("%s", "All params are empty or the same than the actual job")
		return &api.HandlerRes{
			Payload:    err.Error(),
			HttpStatus: 400,
			Err:        err,
		}
	}

	// Stop running job for updating the running cron
	fmt.Println("status: ", job.Status)
	if job.Status == provider.StatusRunning {
		fmt.Println("Status running so stopping...")
		ctx.Manager.Stop(job.ID)
	}

	fmt.Println("Setting up ...")
	err = ctx.Manager.Setup(job)
	if err != nil {
		// This is for keep running the before cronjob that was ok, since it wasn't updated
		ctx.Manager.Start(job.ID)
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}
	fmt.Println("Setted up!")

	fmt.Println("Updating...")
	job.UpdatedAt = time.Now()
	job.Status = provider.StatusRunning
	job, err = ctx.JobStorage.Update(job)
	if err != nil {
		return &api.HandlerRes{Payload: err.Error(), HttpStatus: 500, Err: err}
	}
	fmt.Println("Updated!")

	fmt.Println("Starting ...")
	ctx.Manager.Start(job.ID)
	fmt.Println("Started!")

	res := map[string]interface{}{"id": job.ID, "status": job.Status}
	return &api.HandlerRes{Payload: res, HttpStatus: 200, Err: nil}
}
