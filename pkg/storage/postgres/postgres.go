package postgres

import (
	"Marketplace/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
)

type Storage struct {
	DB  *sql.DB
	log *slog.Logger
}

func NewStorage(conf *config.Config, log *slog.Logger) (*Storage, error) {
	const op = "storage.postgres"

	conStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		conf.PGHost, conf.PGPort, conf.PGUser, conf.PGName, "disable", conf.PGPassword)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}

	return &Storage{DB: db, log: log}, nil
}
