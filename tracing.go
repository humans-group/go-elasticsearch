package es

import (
	"net/http"

	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const (
	errLogKeyEvent   = "event"
	errLogKeyMessage = "message"
	errLogValueErr   = "error"
)

type EsTransportWithTracing struct {
	Name        string
	EsTransport estransport.Interface
}

func (t EsTransportWithTracing) Perform(r *http.Request) (*http.Response, error) {

	ctx := r.Context()

	// operation will include index name, document type and operation itself
	// e.g. bookdb_index/book/_search
	operation := r.URL.Path
	span, spanCtx := opentracing.StartSpanFromContext(ctx, operation)

	spanKind := ext.SpanKindEnum(t.Name) // supposed to be something like "es_client"
	ext.SpanKind.Set(span, spanKind)
	ext.HTTPMethod.Set(span, r.Method)
	defer span.Finish()
	defer spanCtx.Done()

	resp, err := t.EsTransport.Perform(r)

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
