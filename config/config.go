package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PGName         string
	PGPassword     string
	PGUser         string
	PGHost         string
	PGPort         string
	HTTPServerPort string
	RedisDB        string
	RedisPassword  string
	RedisPort      string
	RedisDuration  string
}

func MustLoad() *Config {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	return &Config{
		PGName:         os.Getenv("PG_DB_NAME"),
		PGPassword:     os.Getenv("PG_PASSWORD"),
		PGUser:         os.Getenv("PG_USER"),
		PGHost:         os.Getenv("PG_HOST"),
		PGPort:         os.Getenv("PG_PORT"),
		HTTPServerPort: os.Getenv("HTTP_SERVER_PORT"),
		RedisDB:        os.Getenv("REDIS_DB"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		RedisPort:      os.Getenv("REDIS_PORT"),
		RedisDuration:  os.Getenv("REDIS_DURATION"),
	}
}
