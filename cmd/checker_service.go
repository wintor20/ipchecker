package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"ipchecker/handlers"
	"ipchecker/models"
	"ipchecker/util"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// DBConnectAttempts количество повторов попыток соединения с бд
	DBConnectAttempts = 5
	// DBConnectAttemptInterval интервал между попытками в секундах
	DBConnectAttemptInterval = 5
)

// WaitFor - повторяет вызов ф-ии до победного
// используется для ожидания доступности внешних сервисов
func WaitFor(f func() error, tryInterval int, maxCount int, logger *log.Logger) error {
	for i := 1; i <= maxCount; i++ {
		err := f()
		if err == nil {
			return nil
		}
		logger.Printf("%d попытка соединения с postgres не удалась: %s", i, err.Error())
		time.Sleep(time.Duration(tryInterval) * time.Second)
	}
	return errors.New("превышено число попыток соединения с postgres")
}

// DbConnect пробует соединится с БД, если получается - возвращает объект коннекта, если нет - ошибку
func DbConnect(connStr string, logger *log.Logger) (*sql.DB, error) {
	var db *sql.DB

	err := WaitFor(func() error {
		var err error
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			logger.Printf("%s db connect error %s", connStr, err.Error())
			return err
		}
		err = db.Ping()
		return err
	},
		DBConnectAttemptInterval,
		DBConnectAttempts,
		logger)

	if err != nil {
		return db, err
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxOpenConns(20)
	return db, nil
}

func initService() (*models.ServiceInstance, error) {

	var instance models.ServiceInstance

	var buf bytes.Buffer
	instance.Log = log.New(&buf, "instance log:", log.Lshortfile)
	instance.Log.SetOutput(os.Stdout)

	//------------------------------cfg
	var cfg models.Config

	if err := envconfig.Process("", &cfg); err != nil {
		instance.Log.Fatal(err)
		return nil, err
	}

	instance.HTTPAddr = cfg.HTTPAddr
	instance.HTTPPort = cfg.HTTPPort
	instance.Log.Printf("Config loaded %s", cfg)

	//------------------------------prometheus

	var prometheusMetrics models.Metrics

	prometheusMetrics.Uptime = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "checker_uptime",
		Help: "How many seconds checker worked",
	})

	prometheusMetrics.DeliveredCommands = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "checker_delivered_commands",
		Help: "Total number of delivered commands to checker general channel",
	})

	//-------------------------------LABELED METRICS----------------------------------------------

	prometheusMetrics.FuncUsed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "checker_func_used",
		Help: "Number of func usages",
	}, []string{"function", "code"})

	prometheusMetrics.FuncTimeSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "checker_request_time",
		Help: "Time of request",
	}, []string{"function", "success"})

	//-----------------------------------------------------------------------

	prometheus.MustRegister(prometheusMetrics.Uptime)
	prometheus.MustRegister(prometheusMetrics.FuncUsed)
	prometheus.MustRegister(prometheusMetrics.FuncTimeSummary)

	//-------------------general struct initialize

	instance.PMetrics = &prometheusMetrics

	//-----------------DB stuff
	db, err := DbConnect(cfg.Postgres, instance.Log)
	if err != nil {
		instance.Log.Fatal("не удалось установить соединение с postgres", err)
		return nil, err
	}
	instance.DB = db
	instance.Log.Print("Postgres connected")

	return &instance, nil
}

func main() {

	instance, err := initService()
	if err != nil {
		return
	}

	ipLib, err := util.LoadIPMap(instance)
	if err != nil {
		instance.Log.Fatal(err.Error())
		return
	}

	rout := handlers.DBConnect{DB: instance.DB, IpLib: ipLib}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/isok", rout.IsOk)
	http.HandleFunc("/", rout.CheckIP)

	instance.Log.Print("rest handlers prepared")

	addr := fmt.Sprintf("%s:%s", instance.HTTPAddr, instance.HTTPPort)
	instance.Log.Printf("server prepared on %s", addr)
	srv := &http.Server{Addr: addr}

	instance.Log.Fatal(srv.ListenAndServe())

	return

}
