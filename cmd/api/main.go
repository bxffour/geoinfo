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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.6.1"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

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
	tracer trace.Tracer
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

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", false, "Enable rate limiter")

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

	shutdown, err := initProvider()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer shutdown()

	tracer := otel.GetTracerProvider().Tracer("github.com/bxffour/crest-countries")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		tracer: tracer,
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

func initProvider() (func(), error) {
	ctx := context.Background()

	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("crest-api"),
		),
	)

	if err != nil {
		return nil, err
	}

	otelAddr, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if !ok {
		otelAddr = "0.0.0.0:4317"
	}

	metricExp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(otelAddr),
	)

	if err != nil {
		return nil, err
	}

	MeterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExp,
				sdkmetric.WithInterval(2*time.Second),
			),
		),
	)

	global.SetMeterProvider(MeterProvider)

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelAddr),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)

	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		if err := traceExp.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}

		if err := MeterProvider.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}, nil
}