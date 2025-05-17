package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	DBUrl        string
	GeminiAPIKey string
	JWTSecret    string
}

func Load() *Config {
	c := &Config{
		DBUrl:        os.Getenv("DATABASE_URL"),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
	}
	if c.DBUrl == "" || c.GeminiAPIKey == "" || c.JWTSecret == "" {
		log.Fatal("missing DATABASE_URL, GEMINI_API_KEY or JWT_SECRET")
	}
	return c
}

func ConnectDB(url string) *sql.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	return db
}
