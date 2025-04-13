package model

import "time"

type Item struct {
	ID         string    `db:"id" json:"id"`
	ReceivedAt time.Time `db:"received_at" json:"received_at"`
	Type       string    `db:"type" json:"type"`
}
