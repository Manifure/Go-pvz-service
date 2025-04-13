package model

import "time"

type PVZWithAcceptances struct {
	ID           string       `json:"id"`
	City         string       `json:"city"`
	RegisteredAt time.Time    `db:"registered_at" json:"registered_at"`
	Acceptances  []Acceptance `json:"acceptances"`
}

type PVZQueryParams struct {
	From   string
	To     string
	Limit  int
	Offset int
}
