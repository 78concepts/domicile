package database

import (
	"78concepts.com/domicile/internal/config"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

type DB struct {
	Postgres *pgxpool.Pool
}

func NewPGXPool() *pgxpool.Pool {

	config := config.GetConfig()

	dbPool, err := pgxpool.Connect(context.Background(), config.Database.ConnectionString)

	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	log.Println("Connection to database successful")

	return dbPool
}

