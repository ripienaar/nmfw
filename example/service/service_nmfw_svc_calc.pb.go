// Code generated using Nats Micro Service Framework version 10333f4459fab1d0f306e802b6a126a035fd2fed

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
type AverageHandler func(context.Context, AverageRequest) (*CalcResponse, error)

// Calculates the sum of a series of numbers
type AddHandler func(context.Context, AddRequest) (*CalcResponse, error)

// Compiles and executes a expr expression, expression must return a number
type ExpressionHandler func(context.Context, ExpressionRequest) (*CalcResponse, error)

type CalcService struct {
	Name        string
	Version     string
	Description string
	RootSubject string
	Log         *logrus.Entry

	AverageHandler    AverageHandler
	AddHandler        AddHandler
	ExpressionHandler ExpressionHandler
	ctx               context.Context
	nc                *nats.Conn
}

func (s *CalcService) Start(ctx context.Context, nc *nats.Conn) error {
	wg := sync.WaitGroup{}

	s.ctx = ctx
	s.nc = nc

	wg.Add(1)
	go s.start(&wg, nc)

	wg.Wait()

	return nil
}

func (s *CalcService) handleAverage(req micro.Request) {
	obs := prometheus.NewTimer(handlerRuntime.WithLabelValues(s.Name, "Average"))
	defer obs.ObserveDuration()

	log := s.Log.WithField("function", "Average")
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

	deadline := time.Now().Add(5 * time.Second)
	ds := req.Headers().Get("Nmfw-Deadline")
	if ds != "" {
		deadline, err = time.Parse(time.RFC3339, ds)
		if err != nil {
			log.Errorf("Invalid deadline in request: %v", err)
		}
	}
	log.Debugf("Allowing %v for call to handler", time.Until(deadline))
	to, cancel := context.WithDeadline(s.ctx, deadline)
	defer cancel()

	client := req.Headers().Get("Nmfw-Client-Version")
	rh := &requestContext{
		log: log.WithField("client", client),
		nc:  s.nc,
		req: req,
		cv:  client,
	}

	resp, err := s.AverageHandler(context.WithValue(to, "nmfw", rh), pr)
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

	err = req.Respond(out)
	if err != nil {
		log.Errorf("Publishing response failed: %v", err)
		handlerErrorsCtr.WithLabelValues(s.Name, "Average").Inc()
		return
	}
}

func (s *CalcService) handleAdd(req micro.Request) {
	obs := prometheus.NewTimer(handlerRuntime.WithLabelValues(s.Name, "Add"))
	defer obs.ObserveDuration()

	log := s.Log.WithField("function", "Add")
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

	deadline := time.Now().Add(5 * time.Second)
	ds := req.Headers().Get("Nmfw-Deadline")
	if ds != "" {
		deadline, err = time.Parse(time.RFC3339, ds)
		if err != nil {
			log.Errorf("Invalid deadline in request: %v", err)
		}
	}
	log.Debugf("Allowing %v for call to handler", time.Until(deadline))
	to, cancel := context.WithDeadline(s.ctx, deadline)
	defer cancel()

	client := req.Headers().Get("Nmfw-Client-Version")
	rh := &requestContext{
		log: log.WithField("client", client),
		nc:  s.nc,
		req: req,
		cv:  client,
	}

	resp, err := s.AddHandler(context.WithValue(to, "nmfw", rh), pr)
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

	err = req.Respond(out)
	if err != nil {
		log.Errorf("Publishing response failed: %v", err)
		handlerErrorsCtr.WithLabelValues(s.Name, "Add").Inc()
		return
	}
}

func (s *CalcService) handleExpression(req micro.Request) {
	obs := prometheus.NewTimer(handlerRuntime.WithLabelValues(s.Name, "Expression"))
	defer obs.ObserveDuration()

	log := s.Log.WithField("function", "Expression")
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

	deadline := time.Now().Add(5 * time.Second)
	ds := req.Headers().Get("Nmfw-Deadline")
	if ds != "" {
		deadline, err = time.Parse(time.RFC3339, ds)
		if err != nil {
			log.Errorf("Invalid deadline in request: %v", err)
		}
	}
	log.Debugf("Allowing %v for call to handler", time.Until(deadline))
	to, cancel := context.WithDeadline(s.ctx, deadline)
	defer cancel()

	client := req.Headers().Get("Nmfw-Client-Version")
	rh := &requestContext{
		log: log.WithField("client", client),
		nc:  s.nc,
		req: req,
		cv:  client,
	}

	resp, err := s.ExpressionHandler(context.WithValue(to, "nmfw", rh), pr)
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

	err = req.Respond(out)
	if err != nil {
		log.Errorf("Publishing response failed: %v", err)
		handlerErrorsCtr.WithLabelValues(s.Name, "Expression").Inc()
		return
	}
}

func (s *CalcService) start(wg *sync.WaitGroup, nc *nats.Conn) {
	defer wg.Done()

	svc, err := micro.AddService(nc, micro.Config{
		Name:        s.Name,
		Version:     s.Version,
		Description: s.Description,
		ErrorHandler: func(service micro.Service, natsError *micro.NATSError) {
			s.Log.Errorf("Encountered an error: %v", natsError)
			errorsCtr.WithLabelValues(s.Name, "<no value>").Inc()
		},
	})
	if err != nil {
		panic(fmt.Sprintf("could not create service: %v", err))
	}

	sg := svc.AddGroup("nmfw").AddGroup(s.Name)

	err = sg.AddEndpoint("Average", micro.HandlerFunc(s.handleAverage))
	if err != nil {
		panic(fmt.Sprintf("could not create endpoint Average: %v", err))
	}

	err = sg.AddEndpoint("Add", micro.HandlerFunc(s.handleAdd))
	if err != nil {
		panic(fmt.Sprintf("could not create endpoint Add: %v", err))
	}

	err = sg.AddEndpoint("Expression", micro.HandlerFunc(s.handleExpression))
	if err != nil {
		panic(fmt.Sprintf("could not create endpoint Expression: %v", err))
	}

	nfo := svc.Info()
	s.Log.Infof("Started on subject %v with ID %s", nfo.Subjects, nfo.ID)

	<-s.ctx.Done()
	s.Log.Infof("Shutting down on context")

	svc.Stop()
}

// CalcClient is a client to the CalcService accessed over NATS
type CalcClient struct {
	conn    *nats.Conn
	timeout time.Duration
}

// NewCalcClient creates a new CalcService client using the supplied NATS Connection
func NewCalcClient(nc *nats.Conn, defaultTimeout time.Duration) *CalcClient {
	return &CalcClient{nc, defaultTimeout}
}

func (c *CalcClient) Average(ctx context.Context, req AverageRequest) (*CalcResponse, error) {
	rb, err := proto.Marshal(&req)
	if err != nil {
		return nil, err
	}

	// default to client default timeout but let the context override
	deadline, ok := ctx.Deadline()
	if !ok || deadline.IsZero() {
		deadline = time.Now().Add(c.timeout)
	}
	to, cancel := context.WithTimeout(ctx, time.Until(deadline))
	defer cancel()

	msg := nats.NewMsg("nmfw.calc.Average")
	msg.Data = rb
	msg.Header.Add("Nmfw-Deadline", deadline.Format(time.RFC3339))
	msg.Header.Add("Nmfw-Version", "10333f4459fab1d0f306e802b6a126a035fd2fed")
	msg.Header.Add("Nmfw-Client-Version", "0.0.2")

	srvResp, err := c.conn.RequestMsgWithContext(to, msg)
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

	// default to client default timeout but let the context override
	deadline, ok := ctx.Deadline()
	if !ok || deadline.IsZero() {
		deadline = time.Now().Add(c.timeout)
	}
	to, cancel := context.WithTimeout(ctx, time.Until(deadline))
	defer cancel()

	msg := nats.NewMsg("nmfw.calc.Add")
	msg.Data = rb
	msg.Header.Add("Nmfw-Deadline", deadline.Format(time.RFC3339))
	msg.Header.Add("Nmfw-Version", "10333f4459fab1d0f306e802b6a126a035fd2fed")
	msg.Header.Add("Nmfw-Client-Version", "0.0.2")

	srvResp, err := c.conn.RequestMsgWithContext(to, msg)
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

	// default to client default timeout but let the context override
	deadline, ok := ctx.Deadline()
	if !ok || deadline.IsZero() {
		deadline = time.Now().Add(c.timeout)
	}
	to, cancel := context.WithTimeout(ctx, time.Until(deadline))
	defer cancel()

	msg := nats.NewMsg("nmfw.calc.Expression")
	msg.Data = rb
	msg.Header.Add("Nmfw-Deadline", deadline.Format(time.RFC3339))
	msg.Header.Add("Nmfw-Version", "10333f4459fab1d0f306e802b6a126a035fd2fed")
	msg.Header.Add("Nmfw-Client-Version", "0.0.2")

	srvResp, err := c.conn.RequestMsgWithContext(to, msg)
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
