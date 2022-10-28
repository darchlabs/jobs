package job

import "time"

type Job struct {
	Id           string    `json:"id"`
	ProviderId   string    `json:"providerId"`
	Status       string    `json:"status"`
	Network      string    `json:"network"`
	Address      string    `json:"address"`
	Abi          string    `json:"abi"`
	Type         string    `json:"type,omitempty"`
	Cronjob      string    `json:"cronJob,omitempty"`
	CheckMethod  string    `json:"checkMethod"`
	ActionMethod string    `json:"actionMethod"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
}
