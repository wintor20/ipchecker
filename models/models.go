package models

import (
	"database/sql"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Postgres string `default:""`
	HTTPAddr string `default:""`
	HTTPPort string `default:""`
}

type Metrics struct {
	Uptime prometheus.Counter

	DeliveredCommands prometheus.Gauge

	FuncUsed *prometheus.CounterVec

	FuncTimeSummary *prometheus.SummaryVec
}

type ServiceInstance struct {
	DB *sql.DB

	PMetrics *Metrics
	Log      *log.Logger

	HTTPAddr string
	HTTPPort string
}
