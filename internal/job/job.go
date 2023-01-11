package job

import (
	"time"

	"github.com/darchlabs/jobs/internal/provider"
)

type Job struct {
	ID           string         `json:"id,omitempty"`
	Name         string         `json:"name"`
	ProviderId   string         `json:"providerId"`
	Status       provider.State `json:"status,omitempty"`
	Network      string         `json:"network"`
	Address      string         `json:"address"`
	Abi          string         `json:"abi"`
	NodeURL      string         `json:"nodeUrl"`
	Privatekey   string         `json:"privateKey"`
	Type         string         `json:"type"`
	Cronjob      string         `json:"cronjob,omitempty"`
	CheckMethod  *string        `json:"checkMethod"`
	ActionMethod string         `json:"actionMethod"`
	CreatedAt    time.Time      `json:"createdAt,omitempty"`
	UpdatedAt    time.Time      `json:"updatedAt,omitempty"`
	Logs         []string       `json:"logs,omitempty"`
}
