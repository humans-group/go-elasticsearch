package es

import (
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"

	"github.com/prometheus/client_golang/prometheus"
)

var clientDurationSummary *prometheus.SummaryVec

func init() {
	clientDurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "search_operations_durations_seconds",
			Help: "search_operations_durations_seconds",
		},
		[]string{"client_name", "method"},
	)

	prometheus.MustRegister(clientDurationSummary)
}

func newSearchFuncWithMetrics(searchFunc esapi.Search) esapi.Search {
	return func(o ...func(*esapi.SearchRequest)) (*esapi.Response, error) {
		start := time.Now()

		r, err := searchFunc(o...)

		duration := time.Since(start)
		clientDurationSummary.WithLabelValues("ES", "search").Observe(duration.Seconds())

		return r, err
	}
}

//func (ta *metricsAdapter) observe(method string, startedAt time.Time) {
//	duration := time.Since(startedAt)
//	clientDurationSummary.WithLabelValues(ta.name, method).Observe(duration.Seconds())
//}

// prometheusCollector exports metrics from db.DBStats as prometheus` gauges.
type prometheusCollector struct {
	mu             sync.RWMutex
	searchRequests *prometheus.Desc
}

//var errAlreadyRegistered = errors.New("already registered")

// register adds connection to pool. Returns an error on duplicate pool name.
//func (pc *prometheusCollector) register(name string, conn *pgxpool.Pool) error {
//	if name == "" {
//		name = "default"
//	}
//
//	pc.mu.Lock()
//	defer pc.mu.Unlock()
//
//	if _, exists := pc.dbs[name]; exists {
//		return errAlreadyRegistered
//	}
//
//	pc.dbs[name] = conn
//
//	return nil
//}

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
