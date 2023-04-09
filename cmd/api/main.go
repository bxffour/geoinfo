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
	"github.com/spf13/viper"

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

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	wg     sync.WaitGroup
	tracer trace.Tracer
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	var cfg config

	if err := setupFlags(); err != nil {
		logger.PrintFatal(err, nil)
	}

	if err := setupConfig(&cfg); err != nil {
		logger.PrintFatal(err, nil)
	}

	if viper.GetBool("version") {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		fmt.Printf("Go Version:\t%s\n", goVersion)
		os.Exit(0)
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	logger.PrintInfo("database connection pool established", nil)

	if viper.GetBool("telemetry") {
		shutdown, err := initProvider()
		if err != nil {
			logger.PrintFatal(err, nil)
		}

		defer shutdown()
	}

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
	pgInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s sslcert=%s sslkey=%s",
		cfg.db.host, cfg.db.port, cfg.db.user, cfg.db.password, cfg.db.dbname, cfg.db.sslmode, cfg.db.sslrootcert, cfg.db.sslcert, cfg.db.sslkey)

	db, err := sql.Open("postgres", pgInfo)
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