package metrics

import "github.com/prometheus/client_golang/prometheus"

var ResponseTime = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Name: "response_time",
	Help: "Time request executed",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"endpoint"},)

var Errors = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "errors",
}, []string{"path", "error"})

var Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "hits",
}, []string{"path", "status"})

