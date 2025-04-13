package handler

import (
	"Go-pvz-service/internal/db"
	"Go-pvz-service/internal/model"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"message": invalid request}`, http.StatusBadRequest)
		return
	}

	input.Role = strings.ToLower(input.Role)
	if input.Role != "client" && input.Role != "moderator" && input.Role != "employee" {
		http.Error(w, `{"message": invalid role}`, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"message": server error}`, http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	_, err = db.DB.Exec(
		"INSERT INTO users (id, email, password, role) VALUES ($1, $2, $3, $4)",
		id, input.Email, string(hashedPassword), input.Role,
	)
	if err != nil {
		if isUniqueViolation(err) {
			http.Error(w, `{"message": email already exists"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"message": failed to create user}`+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := model.User{
		ID:    id,
		Email: input.Email,
		Role:  input.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, `{"message": server error}`, http.StatusInternalServerError)
		return
	}

}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "unique")
}
