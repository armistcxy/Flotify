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

func GetAuthDatabasePool() *pgxpool.Pool {
	auth_database_config := config.LoadAuthDatabaseConfig()
	authdbpool, err := pgxpool.New(context.Background(), auth_database_config.DSN)

	if err != nil {
		panic(err)
	}
	err = authdbpool.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	return authdbpool
}

// It seem like bad to have two seperate GetDatabasePool at first glance
// But each config is different so you can ignore it ...
