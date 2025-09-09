package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() {
	url := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	DB = pool
}
