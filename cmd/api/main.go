package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bxffour/crest-countries/internal/data"
	"github.com/bxffour/crest-countries/internal/jsonlog"
	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"

	_ "github.com/lib/pq"
)

var (
	buildTime string
	version   string
	goVersion string
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	var dsn string

	flag.IntVar(&cfg.port, "port", 8080, "server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment ( development | staging | production )")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "Postgres database connection string")
	flag.StringVar(&dsn, "dsn-path", "/etc/crest/.env", "path to .env file containing database connection string")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "Postgresql maximum open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "Postgresql maximum idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "Postgresql maximum idle time")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		fmt.Printf("Go Version:\t%s\n", goVersion)
		os.Exit(0)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	if cfg.db.dsn == "" {
		err := cfg.loadEnv(dsn)
		if err != nil {
			logger.PrintFatal(err, nil)
		}
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c *config) loadEnv(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	dsn := os.Getenv("CRESTCOUNTRIES_DB_DSN")
	c.db.dsn = dsn

	return nil
}
