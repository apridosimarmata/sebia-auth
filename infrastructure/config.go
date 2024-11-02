package infrastructure

import (
	"os"
)

type Config struct {
	POSTGRES_DB       string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string

	REDIS_HOST string
	REDIS_PORT string

	WALLET_TRANSACTION_CHANNEL string
}

func GetConfig() Config {
	return Config{
		POSTGRES_DB:                os.Getenv("POSTGRES_DB"),
		POSTGRES_HOST:              os.Getenv("POSTGRES_HOST"),
		POSTGRES_PORT:              os.Getenv("POSTGRES_PORT"),
		POSTGRES_USER:              os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD:          os.Getenv("POSTGRES_PASSWORD"),
		REDIS_HOST:                 os.Getenv("REDIS_HOST"),
		REDIS_PORT:                 os.Getenv("REDIS_PORT"),
		WALLET_TRANSACTION_CHANNEL: os.Getenv("WALLET_TRANSACTION_CHANNEL"),
	}
}
