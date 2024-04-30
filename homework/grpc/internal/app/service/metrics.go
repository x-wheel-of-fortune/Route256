package service

import "github.com/prometheus/client_golang/prometheus"

var (
	// Create a metrics registry.
	Reg = prometheus.NewRegistry()

	// Create a customized counter metric.
	AddedPointsMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "added_count",
		Help: "Total number of pickup points added.",
	})

	DeletedPointsMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "deleted_count",
		Help: "Total number of pickup points deleted.",
	})

	InternalErrorCountMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "internal_error_count",
		Help: "Total number of server internal errors.",
	})

	ClientErrorCountMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "client_error_count",
		Help: "Total number of errors caused by clients' input.",
	})
)
