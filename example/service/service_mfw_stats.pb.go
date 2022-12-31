// Code generated using Nats Micro Service Framework version <no value>

package service

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	handlerRuntime = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "mwf_handler_runtime",
		Help: "Time the handler took to process",
	}, []string{"service", "function"})

	errorsCtr = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "mwf_errors",
		Help: "Times errors were encountered",
	}, []string{"service", "function"})

	handlerErrorsCtr = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "mwf_handler_errors",
		Help: "Times handler errors were encountered",
	}, []string{"service", "function"})
)

func init() {
	prometheus.MustRegister(handlerRuntime)
	prometheus.MustRegister(handlerErrorsCtr)
	prometheus.MustRegister(errorsCtr)
}