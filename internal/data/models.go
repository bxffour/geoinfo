package data

import (
	"database/sql"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Countries CountryModel
}

func NewModels(db *sql.DB) Models {
	tracer = otel.GetTracerProvider().Tracer("github.com/bxffour/crest-countries/internal/data")

	return Models{
		Countries: CountryModel{DB: db},
	}
}
