package main

import (
	"Go-pvz-service/internal/auth"
	"Go-pvz-service/internal/config"
	"Go-pvz-service/internal/db"
	"Go-pvz-service/internal/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	config.LoadEnv()

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

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
	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
