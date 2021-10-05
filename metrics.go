package es

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7/estransport"
	"github.com/prometheus/client_golang/prometheus"
)

var clientDurationSummary *prometheus.SummaryVec

func init() {
	clientDurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "es_search_operations_durations_seconds",
			Help: "es_search_operations_durations_seconds",
		},
		[]string{"client_name", "operation", "status", "resp_code"},
	)

	prometheus.MustRegister(clientDurationSummary)
}

type EsTransportWithMetrics struct {
	EsTransport estransport.Interface
}

func (t EsTransportWithMetrics) Perform(r *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := t.Perform(r)

	status := "ok"
	if err != nil {
		status = "error"
	}

	operation := "search"
	if !strings.Contains(r.URL.Path, "_search") {
		operation = "other"
	}

	respCode := 0
	if resp != nil {
		respCode = resp.StatusCode
	}

	duration := time.Since(start)
	clientDurationSummary.WithLabelValues(
		"ES", operation, status, strconv.Itoa(respCode)).Observe(duration.Seconds())

	return resp, err
}

// prometheusCollector exports metrics as prometheus gauges.
type prometheusCollector struct {
	mu             sync.RWMutex
	searchRequests *prometheus.Desc
}

var collector *prometheusCollector

func init() {
	collector = &prometheusCollector{
		searchRequests: prometheus.NewDesc("search_requests", "search requests number", []string{"name"}, nil),
	}

	prometheus.MustRegister(collector)
}

// Describe prometheus.Collector interface implementation
func (pc *prometheusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- pc.searchRequests
}

// Collect prometheus.Collector interface implementation
func (pc *prometheusCollector) Collect(ch chan<- prometheus.Metric) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	searchRequestsCount := 100 // ToDo: use real metric value here
	ch <- prometheus.MustNewConstMetric(pc.searchRequests, prometheus.GaugeValue, float64(searchRequestsCount), "Count")
}
