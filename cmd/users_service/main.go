package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/tushar7789/E_commerce_ms/internal/env"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := appConfig{
		addr: env.GetEnv("ADDR", ":8081"),
		db: dbConfig{
			dsn: env.GetEnv("GOOSE_DB_DIR", "./data/db/users.db"),
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := sql.Open("sqlite", cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	logger.Info("connected to database", "dns", cfg.db.dsn)

	api := application{
		config: cfg,
		db:     db,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("server has failed to start", "error", err)
		os.Exit(1)
	}
}
