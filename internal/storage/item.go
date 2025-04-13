package storage

import (
	"Go-pvz-service/internal/model"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

func GetOpenAcceptanceID(db *sqlx.DB, pvzID string) (string, error) {
	var id string
	query := `
		SELECT id FROM acceptances 
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY created_at DESC
		LIMIT 1;
	`
	err := db.Get(&id, query, pvzID)
	return id, err
}

func AddItemToAcceptance(db *sqlx.DB, item model.Item, acceptanceID string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO items (id, received_at, type) 
		VALUES ($1, $2, $3)
	`, item.ID, item.ReceivedAt, item.Type)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO acceptance_items (acceptance_id, item_id) 
		VALUES ($1, $2)
	`, acceptanceID, item.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetLastAcceptanceIDAndStatus(db *sqlx.DB, pvzID string) (string, string, error) {
	var id string
	var status string

	err := db.QueryRow(
		`SELECT id, status FROM acceptances 
		 WHERE pvz_id = $1 
		 ORDER BY created_at DESC 
		 LIMIT 1`, pvzID).Scan(&id, &status)

	if err == sql.ErrNoRows {
		return "", "", nil
	}

	if err != nil {
		return "", "", err
	}

	return id, status, nil
}

func DeleteLastItemFromAcceptance(db *sqlx.DB, acceptanceID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var itemID string
	err = tx.QueryRow(
		`SELECT item_id FROM acceptance_items 
		 JOIN items ON items.id = acceptance_items.item_id
		 WHERE acceptance_id = $1
		 ORDER BY items.received_at DESC
		 LIMIT 1`, acceptanceID).Scan(&itemID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM acceptance_items WHERE acceptance_id = $1 AND item_id = $2`, acceptanceID, itemID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM items WHERE id = $1`, itemID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
