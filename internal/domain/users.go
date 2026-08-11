package domain

import (
	"database/sql"
)

type application struct {
	config config
	//logger
	db *sql.DB
}

type config struct {
	addr string // addr of server
	db   dbConfig
}

type dbConfig struct {
	dsn string // user, pass, domain
}