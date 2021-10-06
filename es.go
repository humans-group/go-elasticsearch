/*
This package is an entry point in the same way as package elasticsearch in the official go client.
Different name is chosen to avoid corresponding warning.
Acquiring client is different only in package name and type of passed Config struct:
	esClient, err := es.NewClient(es.Config{})
*/
package es

import (
	"github.com/elastic/go-elasticsearch/v7"
)

// Returning standard Client from elasticsearch package but with modified transport
func NewClient(serviceEsCfg Config) (*elasticsearch.Client, error) {
	var client *elasticsearch.Client
	var err error

	cfg := elasticsearch.Config{
		Addresses: serviceEsCfg.Addresses,
	}

	client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	if serviceEsCfg.Tracing {
		client.Transport = EsTransportWithTracing{EsTransport: client.Transport}
	}

	if serviceEsCfg.Metrics {
		client.Transport = EsTransportWithMetrics{
			Name:        serviceEsCfg.Name,
			EsTransport: client.Transport}
	}

	return client, nil
}

func MustNew(serviceEsCfg Config) (*elasticsearch.Client, error) {
	client, err := NewClient(serviceEsCfg)
	if err != nil {
		panic(err)
	}
	return client, nil
}
