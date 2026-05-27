package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	RequestTotal    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

func New(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		RequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of http requests",
			},
			[]string{"method", "path", "status"},
		),

		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "http request latecy",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}

	reg.MustRegister(
		m.RequestTotal,
		m.RequestDuration,
	)

	return m
}
