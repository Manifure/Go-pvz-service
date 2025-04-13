package storage

import (
	"Go-pvz-service/internal/model"
	"github.com/jmoiron/sqlx"
)

func GetPVZWithAcceptancesFiltered(db *sqlx.DB, params model.PVZQueryParams) ([]model.PVZ, error) {
	var pvzList []model.PVZ

	query := `SELECT id, city, registered_at FROM pvz ORDER BY registered_at DESC LIMIT $1 OFFSET $2`
	err := db.Select(&pvzList, query, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	for i := range pvzList {
		var acceptances []model.Acceptance
		query := `
			SELECT id, created_at, pvz_id, status
			FROM acceptances
			WHERE pvz_id = $1 AND created_at BETWEEN $2 AND $3
			ORDER BY created_at DESC
		`
		err := db.Select(&acceptances, query, pvzList[i].ID, params.From, params.To)
		if err != nil {
			return nil, err
		}
		pvzList[i].Acceptances = acceptances
	}

	return pvzList, nil
}
