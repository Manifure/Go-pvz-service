package model

import "time"

type PVZ struct {
	ID           string       `db:"id" json:"id"`
	City         string       `db:"city" json:"city"`
	RegisteredAt time.Time    `db:"registered_at" json:"registered_at"`
	Acceptances  []Acceptance `json:"acceptances"`
}
