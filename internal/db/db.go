package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sqlx.DB

func Init() error {
	var err error
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("POSTGRES_HOST", "postgres"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_USER", "postgres"),
		getEnv("POSTGRES_PASSWORD", "1"),
		getEnv("POSTGRES_NAME", "postgres"),
	)
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
		return err
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
		return err
	}
	fmt.Println("Successfully connected to DB")
	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
