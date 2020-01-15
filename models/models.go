package models

import (
	"database/sql"
	"log"
)

type Config struct {
	Postgres string `default:"postgres://checker:checker@localhost/checker_db?sslmode=disable"`
	HTTPAddr string `default:"localhost"`
	HTTPPort string `default:"8098"`
}

type ServiceInstance struct {
	DB       *sql.DB
	Log      *log.Logger
	HTTPAddr string
	HTTPPort string
}

type CheckIPAnswer struct {
	Dupes bool `json:"dupes"`
}
