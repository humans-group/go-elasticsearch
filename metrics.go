package es

import (
	"net/http"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	statusOk    = "ok"
	statusError = "error"
)

var clientDurationSummary *prometheus.SummaryVec

func init() {
	clientDurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "es_operations_durations_seconds",
			Help: "es_operations_durations_seconds",
		},
		[]string{"client_name", "operation", "status", "resp_code"},
	)

	prometheus.MustRegister(clientDurationSummary)
}

type EsTransportWithMetrics struct {
	Name        string
	EsTransport estransport.Interface
}

func (t EsTransportWithMetrics) Perform(r *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := t.Perform(r)

	status := statusOk
	if err != nil {
		status = statusError
	}

	// operation will include index name, document type and operation itself
	// e.g. bookdb_index/book/_search
	operation := r.URL.Path

	respCode := 0
	if resp != nil {
		respCode = resp.StatusCode
	}

	duration := time.Since(start)
	clientDurationSummary.WithLabelValues(
		t.Name, operation, status, strconv.Itoa(respCode)).Observe(duration.Seconds())

	return resp, err
}
