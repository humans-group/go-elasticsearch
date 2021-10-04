package es

import (
	"context"

	"github.com/elastic/go-elasticsearch/v7/esapi"
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

func newSearchFuncWithTracing(searchFunc esapi.Search) esapi.Search {
	return func(o ...func(*esapi.SearchRequest)) (*esapi.Response, error) {

		// todo: get context from prepared request?
		ctx := context.Background()

		span, _ := startSpan(ctx, operationNameSearch)
		r, err := searchFunc(o...)

		if err != nil {
			traceErr(err, span)
		}

		span.Finish()

		return r, err
	}
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
