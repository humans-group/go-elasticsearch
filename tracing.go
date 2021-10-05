package es

import (
	"context"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const (
	operationNameSearch = "es.Search"
	errLogKeyEvent      = "event"
	errLogKeyMessage    = "message"
	errLogValueErr      = "error"
)

type EsTransportWithTracing struct {
	EsTransport estransport.Interface
}

func (t EsTransportWithTracing) Perform(r *http.Request) (*http.Response, error) {

	ctx := r.Context()
	span, spanCtx := startSpan(ctx, operationNameSearch)
	defer span.Finish()
	defer spanCtx.Done()

	resp, err := t.Perform(r)

	if err != nil {
		traceErr(err, span)
	}

	return resp, err
}

func traceErr(err error, span opentracing.Span) {
	ext.Error.Set(span, true)
	span.LogFields(
		log.String(errLogKeyEvent, errLogValueErr),
		log.String(errLogKeyMessage, err.Error()),
	)
}

func startSpan(ctx context.Context, name string) (opentracing.Span, context.Context) {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, name)
	return span, spanCtx
}
