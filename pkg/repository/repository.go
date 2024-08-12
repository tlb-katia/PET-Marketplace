package repository

import (
	"Marketplace/pkg/storage/postgres"
	"database/sql"
	"github.com/go-redis/redis/v8"
)

type Repository struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewRepository(storage *postgres.Storage, rdb *redis.Client) *Repository {
	return &Repository{
		db:  storage.DB,
		rdb: rdb,
	}
}
