package model

import "time"

type Acceptance struct {
	ID        string    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	PVZID     string    `db:"pvz_id" json:"pvz_id"`
	Status    string    `db:"status" json:"status"`
	Items     []Item    `json:"items"`
}

type AcceptanceWithItems struct {
	Acceptance
	Items []Item `json:"items"`
}
