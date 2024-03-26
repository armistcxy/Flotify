package database

import (
	"context"
	"flotify/internal/config"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDatabasePool() *pgxpool.Pool {
	database_config := config.LoadDatabaseConfig()
	dbpool, err := pgxpool.New(context.Background(), database_config.DSN)
	if err != nil {
		panic(err)
	}
	err = dbpool.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	log.Println("Connect database successfully")
	return dbpool
}
