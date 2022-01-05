package main

import (
	"flag"
	"os"
	"sync"

	"github.com/Nana-Seyram/crest-countries/internal/jsonlog"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment ( development | staging | production )")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
