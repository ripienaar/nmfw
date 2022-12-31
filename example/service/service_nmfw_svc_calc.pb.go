// Code generated using Nats Micro Service Framework version development

package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// Calculates the average of a series of numbers
type AverageHandler func(AverageRequest) (*CalcResponse, error)

// Calculates the sum of a series of numbers
type AddHandler func(AddRequest) (*CalcResponse, error)

// Compiles and executes a expr expression, expression must return a number
type ExpressionHandler func(ExpressionRequest) (*CalcResponse, error)

type CalcService struct {
	Name        string
	Version     string
	Description string
	RootSubject string
	Log         *logrus.Entry

	AverageHandler    AverageHandler
	AddHandler        AddHandler
	ExpressionHandler ExpressionHandler
}

func (s *CalcService) Start(ctx context.Context, nc *nats.Conn) error {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go s.startAverageHandler(ctx, &wg, nc)

	wg.Add(1)
	go s.startAddHandler(ctx, &wg, nc)

	wg.Add(1)
	go s.startExpressionHandler(ctx, &wg, nc)

	wg.Wait()

	return nil
}

func (s *CalcService) startAverageHandler(ctx context.Context, wg *sync.WaitGroup, nc *nats.Conn) {
	defer wg.Done()

	name := fmt.Sprint(s.Name, "_Average")
	timer := handlerRuntime.WithLabelValues(s.Name, "Average")
	log := s.Log.WithField("function", "Average")

	svc, err := micro.AddService(nc, micro.Config{
		Name:        name,
		Version:     s.Version,
		Description: s.Description,
		ErrorHandler: func(service micro.Service, natsError *micro.NATSError) {
			log.Errorf("Encountered an error: %v", natsError)
			errorsCtr.WithLabelValues(s.Name, "Average").Inc()
		},
		Endpoint: micro.Endpoint{
			Subject: fmt.Sprintf("%s.Average", s.RootSubject),
			Handler: func(req *micro.Request) {
				obs := prometheus.NewTimer(timer)
				defer obs.ObserveDuration()

				log.Infof("Handling a request")

				if s.AverageHandler == nil {
					handlerErrorsCtr.WithLabelValues(s.Name, "Average").Inc()
					req.Error("NO_IMPL", "not implemented", nil)
					return
				}

				var pr AverageRequest
				err := proto.Unmarshal(req.Data(), &pr)
				if err != nil {
					log.Errorf("Unmarshaling request failed: %v", err)
					handlerErrorsCtr.WithLabelValues(s.Name, "Average").Inc()
					req.Error("REQ_INVALID", "invalid request", nil)
					return
				}

				resp, err := s.AverageHandler(pr)
				if err != nil {
					log.Errorf("Handling request failed: %v", err)
					handlerErrorsCtr.WithLabelValues(s.Name, "Average").Inc()
					req.Error("ERROR", err.Error(), nil)
					return
				}

				out, err := proto.Marshal(resp)
				if err != nil {
					log.Errorf("Marshaling request failed: %v", err)
					handlerErrorsCtr.WithLabelValues(s.Name, "Average").Inc()
					req.Error("ERROR", err.Error(), nil)
					return
				}

				req.Respond(out)
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Average handler error: %v", err))
	}

	nfo := svc.Info()
	log.Infof("Started on subject %s with ID %s", nfo.Subject, nfo.ID)

	<-ctx.Done()
	log.Infof("Shutting down on context")

	svc.Stop()
}

func (s *CalcService) startAddHandler(ctx context.Context, wg *sync.WaitGroup, nc *nats.Conn) {
	defer wg.Done()

	name := fmt.Sprint(s.Name, "_Add")
	timer := handlerRuntime.WithLabelValues(s.Name, "Add")
	log := s.Log.WithField("function", "Add")

	svc, err := micro.AddService(nc, micro.Config{
		Name:        name,
		Version:     s.Version,
		Description: s.Description,
		ErrorHandler: func(service micro.Service, natsError *micro.NATSError) {
			log.Errorf("Encountered an error: %v", natsError)
			errorsCtr.WithLabelValues(s.Name, "Add").Inc()
		},
		Endpoint: micro.Endpoint{
			Subject: fmt.Sprintf("%s.Add", s.RootSubject),
			Handler: func(req *micro.Request) {
				obs := prometheus.NewTimer(timer)
				defer obs.ObserveDuration()

				log.Infof("Handling a request")

				if s.AddHandler == nil {
					handlerErrorsCtr.WithLabelValues(s.Name, "Add").Inc()
					req.Error("NO_IMPL", "not implemented", nil)
					return
				}

				var pr AddRequest
				err := proto.Unmarshal(req.Data(), &pr)
				if err != nil {
					log.Errorf("Unmarshaling request failed: %v", err)
					handlerErrorsCtr.WithLabelValues(s.Name, "Add").Inc()
					req.Error("REQ_INVALID", "invalid request", nil)
					return
				}

				resp, err := s.AddHandler(pr)
				if err != nil {
					log.Errorf("Handling request failed: %v", err)
					handlerErrorsCtr.WithLabelValues(s.Name, "Add").Inc()
					req.Error("ERROR", err.Error(), nil)
					return
				}

				out, err := proto.Marshal(resp)
				if err != nil {
					log.Errorf("Marshaling request failed: %v", err)
					handlerErrorsCtr.WithLabelValues(s.Name, "Add").Inc()
					req.Error("ERROR", err.Error(), nil)
					return
				}

				req.Respond(out)
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Add handler error: %v", err))
	}

	nfo := svc.Info()
	log.Infof("Started on subject %s with ID %s", nfo.Subject, nfo.ID)

	<-ctx.Done()
	log.Infof("Shutting down on context")

	svc.Stop()
}

func (s *CalcService) startExpressionHandler(ctx context.Context, wg *sync.WaitGroup, nc *nats.Conn) {
	defer wg.Done()

	name := fmt.Sprint(s.Name, "_Expression")
	timer := handlerRuntime.WithLabelValues(s.Name, "Expression")
	log := s.Log.WithField("function", "Expression")

	svc, err := micro.AddService(nc, micro.Config{
		Name:        name,
		Version:     s.Version,
		Description: s.Description,
		ErrorHandler: func(service micro.Service, natsError *micro.NATSError) {
			log.Errorf("Encountered an error: %v", natsError)
			errorsCtr.WithLabelValues(s.Name, "Expression").Inc()
		},
		Endpoint: micro.Endpoint{
			Subject: fmt.Sprintf("%s.Expression", s.RootSubject),
			Handler: func(req *micro.Request) {
				obs := prometheus.NewTimer(timer)
				defer obs.ObserveDuration()

				log.Infof("Handling a request")

				if s.ExpressionHandler == nil {
					handlerErrorsCtr.WithLabelValues(s.Name, "Expression").Inc()
					req.Error("NO_IMPL", "not implemented", nil)
					return
				}

				var pr ExpressionRequest
				err := proto.Unmarshal(req.Data(), &pr)
				if err != nil {
					log.Errorf("Unmarshaling request failed: %v", err)
					handlerErrorsCtr.WithLabelValues(s.Name, "Expression").Inc()
					req.Error("REQ_INVALID", "invalid request", nil)
					return
				}

				resp, err := s.ExpressionHandler(pr)
				if err != nil {
					log.Errorf("Handling request failed: %v", err)
					handlerErrorsCtr.WithLabelValues(s.Name, "Expression").Inc()
					req.Error("ERROR", err.Error(), nil)
					return
				}

				out, err := proto.Marshal(resp)
				if err != nil {
					log.Errorf("Marshaling request failed: %v", err)
					handlerErrorsCtr.WithLabelValues(s.Name, "Expression").Inc()
					req.Error("ERROR", err.Error(), nil)
					return
				}

				req.Respond(out)
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Expression handler error: %v", err))
	}

	nfo := svc.Info()
	log.Infof("Started on subject %s with ID %s", nfo.Subject, nfo.ID)

	<-ctx.Done()
	log.Infof("Shutting down on context")

	svc.Stop()
}

// CalcClient is a client to the CalcService accessed over NATS
type CalcClient struct {
	conn    *nats.Conn
	timeout time.Duration
}

// NewCalcClient creates a new CalcService client using the supplied NATS Connection
func NewCalcClient(nc *nats.Conn, timeout time.Duration) *CalcClient {
	return &CalcClient{nc, timeout}
}

func (c *CalcClient) Average(ctx context.Context, req AverageRequest) (*CalcResponse, error) {
	rb, err := proto.Marshal(&req)
	if err != nil {
		return nil, err
	}

	to, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	srvResp, err := c.conn.RequestWithContext(to, "nmfw.calc.Average", rb)
	if err != nil {
		return nil, err
	}

	if eh := srvResp.Header.Get(micro.ErrorHeader); eh != "" {
		return nil, fmt.Errorf(eh)
	}

	var res CalcResponse
	err = proto.Unmarshal(srvResp.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
func (c *CalcClient) Add(ctx context.Context, req AddRequest) (*CalcResponse, error) {
	rb, err := proto.Marshal(&req)
	if err != nil {
		return nil, err
	}

	to, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	srvResp, err := c.conn.RequestWithContext(to, "nmfw.calc.Add", rb)
	if err != nil {
		return nil, err
	}

	if eh := srvResp.Header.Get(micro.ErrorHeader); eh != "" {
		return nil, fmt.Errorf(eh)
	}

	var res CalcResponse
	err = proto.Unmarshal(srvResp.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
func (c *CalcClient) Expression(ctx context.Context, req ExpressionRequest) (*CalcResponse, error) {
	rb, err := proto.Marshal(&req)
	if err != nil {
		return nil, err
	}

	to, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	srvResp, err := c.conn.RequestWithContext(to, "nmfw.calc.Expression", rb)
	if err != nil {
		return nil, err
	}

	if eh := srvResp.Header.Get(micro.ErrorHeader); eh != "" {
		return nil, fmt.Errorf(eh)
	}

	var res CalcResponse
	err = proto.Unmarshal(srvResp.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}