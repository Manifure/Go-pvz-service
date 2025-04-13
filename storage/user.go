package storage

import (
	"Go-pvz-service/internal/model"
	"github.com/jmoiron/sqlx"
)

func GetUserByEmail(db *sqlx.DB, email string) (*model.User, error) {
	var user model.User
	err := db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	return &user, err
}
