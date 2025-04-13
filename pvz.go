package storage

import (
	"Go-pvz-service/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

var allowedCities = map[string]bool{
	"Москва":          true,
	"Санкт-Петербург": true,
	"Казань":          true,
}

func CreatePVZ(db *sqlx.DB, city string) (*model.PVZ, error) {
	if !allowedCities[city] {
		return nil, ErrCityNotAllowed
	}

	pvz := &model.PVZ{
		ID:           uuid.New().String(),
		City:         city,
		RegisteredAt: time.Now(),
	}

	_, err := db.Exec(
		`INSERT INTO pvz (id, city, registered_at) VALUES ($1, $2, $3)`,
		pvz.ID, pvz.City, pvz.RegisteredAt,
	)
	if err != nil {
		return nil, err
	}
	return pvz, nil
}
