package storage

import (
	"Go-pvz-service/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

func HasOpenAcceptance(db *sqlx.DB, pvzID string) (bool, error) {
	var count int
	err := db.Get(&count,
		`SELECT COUNT(*) FROM acceptances WHERE pvz_id = $1 AND status = 'in_progress'`,
		pvzID,
	)
	return count > 0, err
}

func CreateAcceptanceWithItems(db *sqlx.DB, acceptance *model.Acceptance) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(
		`INSERT INTO acceptances (id, created_at, pvz_id, status) 
		 VALUES (:id, :created_at, :pvz_id, :status)`, acceptance,
	)
	if err != nil {
		return err
	}

	for i := range acceptance.Items {
		item := &acceptance.Items[i]
		item.ID = uuid.New().String()
		item.ReceivedAt = time.Now()

		_, err = tx.Exec(
			`INSERT INTO items (id, received_at, type) VALUES ($1, $2, $3)`,
			item.ID, item.ReceivedAt, item.Type,
		)
		if err != nil {
			return err
		}

		_, err = tx.Exec(
			`INSERT INTO acceptance_items (acceptance_id, item_id) VALUES ($1, $2)`,
			acceptance.ID, item.ID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

//func GetAcceptanceWithItems(db *sqlx.DB, id string) (*model.Acceptance, error) {
//	var acceptance model.Acceptance
//
//	err := db.Get(&acceptance,
//		`SELECT id, created_at, pvz_id, status FROM acceptances WHERE id = $1`,
//		id,
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	err = db.Select(&acceptance.Items,
//		`SELECT i.id, i.received_at, i.type
//		 FROM items i
//		 JOIN acceptance_items ai ON ai.item_id = i.id
//		 WHERE ai.acceptance_id = $1`,
//		id,
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	return &acceptance, nil
//}

func CloseAcceptance(db *sqlx.DB, acceptanceID string) error {
	_, err := db.Exec(`UPDATE acceptances SET status = 'closed' WHERE id = $1`, acceptanceID)
	return err
}
