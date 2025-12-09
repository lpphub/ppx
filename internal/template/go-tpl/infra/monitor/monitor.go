package monitor

import "github.com/prometheus/client_golang/prometheus"

var (
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
)

func Init() {
	prometheus.MustRegister(HTTPRequests)
}