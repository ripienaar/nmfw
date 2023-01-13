// Code generated using Nats Micro Service Framework version 10333f4459fab1d0f306e802b6a126a035fd2fed

package service

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// NatsMicroServicesFrameworkVersion is the version of github.com/ripienaar/nmfw used to generate this service
const NatsMicroServicesFrameworkVersion = "10333f4459fab1d0f306e802b6a126a035fd2fed"

type requestContext struct {
	log *logrus.Entry
	nc  *nats.Conn
	req micro.Request
	cv  string
}

func (r *requestContext) Logger() *logrus.Entry  { return r.log }
func (r *requestContext) Conn() *nats.Conn       { return r.nc }
func (r *requestContext) Request() micro.Request { return r.req }
func (r *requestContext) ClientVersion() string  { return r.cv }

var (
	handlerRuntime = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "nmfw_handler_runtime",
		Help: "Time the handler took to process",
	}, []string{"service", "function"})

	errorsCtr = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "nmfw_errors",
		Help: "Times errors were encountered",
	}, []string{"service", "function"})

	handlerErrorsCtr = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "nmfw_handler_errors",
		Help: "Times handler errors were encountered",
	}, []string{"service", "function"})
)

func init() {
	prometheus.MustRegister(handlerRuntime)
	prometheus.MustRegister(handlerErrorsCtr)
	prometheus.MustRegister(errorsCtr)
}
