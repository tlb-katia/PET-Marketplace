package main

import (
	"Marketplace/config"
	"Marketplace/internal/server"
	"Marketplace/pkg/repository"
	"Marketplace/pkg/storage/postgres"
	"Marketplace/pkg/storage/redis_db"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	conf := config.MustLoad()
	log := setupLogger("local")
	log.Info("starting server", slog.String("address", conf.HTTPServerPort))

	psql, err := postgres.NewStorage(conf, log)
	if err != nil {
		log.Error("failed to init psql", err)
		os.Exit(1)
	}
	defer psql.DB.Close()

	rdb, err := redis_db.NewRedis(conf)
	if err != nil {
		log.Error("failed to init redis", err)
		os.Exit(1)
	}

	repo := repository.NewRepository(psql, rdb.Client)
	srv := server.NewServer(repo, chi.NewRouter(), log)
	srv.Run(conf)

	log.Info("server stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
