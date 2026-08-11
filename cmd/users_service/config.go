package main

import (
	"database/sql"
	"time"

	env "github.com/tushar7789/E_commerce_ms/internal/env"
)

type application struct {
	config appConfig
	//logger
	db *sql.DB
}

type appConfig struct {
	addr string // addr of server
	db   dbConfig
}

type dbConfig struct {
	dsn string // user, pass, domain
}

type tokensConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

func Load() tokensConfig {
	return tokensConfig{
		AccessSecret:  env.GetEnv("ACCESS_SECRET", "dev-access-secret-change-me"),
		RefreshSecret: env.GetEnv("REFRESH_SECRET", "dev-refresh-secret-change-me"),
		AccessTTL:     mustDuration(env.GetEnv("ACCESS_TTL", "15m")),
		RefreshTTL:    mustDuration(env.GetEnv("REFRESH_TTL", "2h")),
	}
}
