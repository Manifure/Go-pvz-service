package handler

import (
	"Go-pvz-service/internal/auth"
	"encoding/json"
	"net/http"
)

func DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Role string `json:"role"`
	}

	var req Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if req.Role != "client" && req.Role != "moderator" && req.Role != "employee" {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateDummyJWT(req.Role)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"token": token}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
