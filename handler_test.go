package test

import (
	"Go-pvz-service/internal/auth"
	"Go-pvz-service/internal/db"
	"Go-pvz-service/internal/handler"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFullAcceptanceFlow(t *testing.T) {
	err := db.Init()
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	router := setupRouter()

	//register(t, router, "moderator@example.com", "password", "moderator")

	moderatorToken := loginAndGetToken(t, router, "moderator@example.com", "123456")

	pvzID := createPVZ(t, router, moderatorToken, "Москва")

	//register(t, router, "employee@example.com", "password", "employee")

	employeeToken := loginAndGetToken(t, router, "employee@example.com", "123456")

	_ = createAcceptance(t, router, employeeToken, pvzID)

	for i := 0; i < 50; i++ {
		addItem(t, router, employeeToken, pvzID, "электроника")
	}

	closeAcceptance(t, router, employeeToken, pvzID)
}

func setupRouter() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/dummyLogin", handler.DummyLoginHandler).Methods("POST")
	r.HandleFunc("/register", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handler.LoginHandler).Methods("POST")
	r.HandleFunc("/pvz", auth.AuthMiddleware(handler.CreatePVZHandler, "moderator")).Methods("POST")
	r.HandleFunc("/acceptances", auth.AuthMiddleware(handler.CreateAcceptanceHandler, "employee")).Methods("POST")
	r.HandleFunc("/items", auth.AuthMiddleware(handler.CreateItemHandler, "employee")).Methods("POST")
	r.HandleFunc("/items", auth.AuthMiddleware(handler.DeleteItemHandler, "employee")).Methods("DELETE")
	r.HandleFunc("/acceptances/close", auth.AuthMiddleware(handler.CloseAcceptanceHandler, "employee")).Methods("POST")
	r.HandleFunc("/info", auth.AuthMiddleware(handler.GetPVZDataHandler, "employee", "moderator")).Methods("GET")

	return r
}

//func register(t *testing.T, router http.Handler, email, password, role string) {
//	body := fmt.Sprintf(`{"email":"%s", "password":"%s", "role":"%s"}`, email, password, role)
//	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
//	req.Header.Set("Content-Type", "application/json")
//
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	require.Equal(t, http.StatusOK, w.Code)
//}

func loginAndGetToken(t *testing.T, router http.Handler, email, password string) string {
	body := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password)
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Token string `json:"token"`
	}
	fmt.Println("Login response body:", w.Body.String())

	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	return resp.Token
}

func createPVZ(t *testing.T, router http.Handler, token, city string) string {
	body := fmt.Sprintf(`{"city":"%s"}`, city)
	req := httptest.NewRequest(http.MethodPost, "/pvz", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		ID string `json:"id"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	return resp.ID
}

func createAcceptance(t *testing.T, router http.Handler, token, pvzID string) string {
	body := fmt.Sprintf(`{"pvz_id":"%s","items":[]}`, pvzID)
	req := httptest.NewRequest(http.MethodPost, "/acceptances", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		ID string `json:"id"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	return resp.ID
}

func addItem(t *testing.T, router http.Handler, token, pvzID, itemType string) {
	body := fmt.Sprintf(`{"pvz_id":"%s", "type":"%s"}`, pvzID, itemType)
	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func closeAcceptance(t *testing.T, router http.Handler, token, pvzID string) {
	body := fmt.Sprintf(`{"pvz_id":"%s"}`, pvzID)
	req := httptest.NewRequest(http.MethodPost, "/acceptances/close", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}
