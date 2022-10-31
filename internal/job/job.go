package job

import "time"

type Job struct {
	ID           string    `json:"id,omitempty"`
	ProviderId   string    `json:"providerId"`
	Status       string    `json:"status"`
	Network      string    `json:"network"`
	Address      string    `json:"address"`
	Abi          string    `json:"abi"`
	Type         string    `json:"type"`
	Cronjob      string    `json:"cronJob,omitempty"`
	CheckMethod  string    `json:"checkMethod"`
	ActionMethod string    `json:"actionMethod"`
	CreatedAt    time.Time `json:"createdAt"`
	UpddatedAt   time.Time `json:"updatedAt,omitempty"`
}
