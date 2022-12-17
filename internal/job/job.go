package job

import "time"

type Job struct {
	ID           string    `json:"id,omitempty"`
	Name         string    `json:"name"`
	ProviderId   string    `json:"providerId"`
	Status       string    `json:"status,omitempty"`
	Network      string    `json:"network"`
	Address      string    `json:"address"`
	Abi          string    `json:"abi"`
	Type         string    `json:"type"`
	Cronjob      string    `json:"cronjob,omitempty"`
	CheckMethod  *string   `json:"checkMethod"`
	ActionMethod string    `json:"actionMethod"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
	Logs         *string   `json:"logs,omitempty"`
}
