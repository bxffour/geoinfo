package data

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var ErrRecordNotFound = errors.New("record not found")

func handleError(ctx context.Context, err error) error {
	span := trace.SpanFromContext(ctx)

	span.RecordError(err, trace.WithTimestamp(time.Now()))
	span.SetStatus(codes.Error, err.Error())

	return err
}

type OtelWrappedError struct {
	err         error
	description string
}

func (ot *OtelWrappedError) Error() string {
	return ot.err.Error()
}

func (ot *OtelWrappedError) Unwrap() error {
	return ot.err
}

func (ot *OtelWrappedError) WithSpanInfo(ctx context.Context, err error) error {
	span := trace.SpanFromContext(ctx)

	span.RecordError(ot.err, trace.WithTimestamp(time.Now()))
	span.SetStatus(codes.Error, ot.description)

	return ot
}