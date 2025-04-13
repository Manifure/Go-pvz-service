package storage_test

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"

	"Go-pvz-service/internal/storage"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("postgres", "host=postgres user=postgres password=1 dbname=testdb sslmode=disable")
	require.NoError(t, err)

	db.Exec("DELETE FROM acceptance_items")
	db.Exec("DELETE FROM items")
	db.Exec("DELETE FROM acceptances")

	return db
}

func TestHasOpenAcceptance(t *testing.T) {
	db := setupTestDB(t)

	pvzID := "pvz-test-id"
	acceptanceID := "acc-test-id"

	_, err := db.Exec(`INSERT INTO acceptances (id, pvz_id, created_at, status) VALUES ($1, $2, $3, $4)`,
		acceptanceID, pvzID, time.Now(), "in_progress")
	require.NoError(t, err)

	open, err := storage.HasOpenAcceptance(db, pvzID)
	require.NoError(t, err)
	require.True(t, open)

	_, err = db.Exec(`UPDATE acceptances SET status = 'in_progress' WHERE id = $1`, acceptanceID)
	require.NoError(t, err)

	open, err = storage.HasOpenAcceptance(db, pvzID)
	require.NoError(t, err)
	require.False(t, open)
}
