package handler

import (
	"Go-pvz-service/internal/auth"
	"Go-pvz-service/internal/db"
	"Go-pvz-service/internal/storage"
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"invalid request"}`, http.StatusBadRequest)
		return
	}

	user, err := storage.GetUserByEmail(db.DB, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"message":"invalid credentials"}`, http.StatusUnauthorized)
		} else {
			http.Error(w, `{"message":"internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, `{"message":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Role)
	if err != nil {
		http.Error(w, `{"message":"failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
