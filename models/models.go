package models

import (
	"database/sql"
	"log"
)

type Config struct {
	Postgres string `default:""`
	HTTPAddr string `default:""`
	HTTPPort string `default:""`
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
