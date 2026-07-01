package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect opens a connection pool to PostgreSQL using the given connection string
// it also verifies if it works or not by pinging the database once before returning.
func Connect(databaseURL string) *pgxpool.Pool {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Never keep more than  10 connections
	// We're broke, we use the neon free tier so we has connection limits
	pool.Config().MaxConns = 10

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	log.Println("Connected to database")
	return pool
}
