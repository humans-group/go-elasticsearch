/*
This package is an entry point in the same way as package elasticsearch in the official go client.
Different name is chosen to avoid corresponding warning.
Acquiring client is different only in package name and type of passed Config struct:
	esClient, err := elasticSearch.NewClient(es.Config{})
*/
package es

import (
	"github.com/elastic/go-elasticsearch/v7"
)

// Returning standard Client from elasticsearch package but with modified Search function
func NewClient(serviceEsCfg Config) (*elasticsearch.Client, error) {
	var client *elasticsearch.Client
	var err error
	client, err = elasticsearch.NewClient(elasticsearch.Config{Addresses: serviceEsCfg.Addresses})

	// Todo: add tracing
	//if serviceEsCfg.Tracing {
	//	client.Search = newSearchFuncWithTracing(client.Search)
	//}

	if serviceEsCfg.Metrics {
		client.Search = newSearchFuncWithMetrics(client.Search)
	}

	return client, err
}
