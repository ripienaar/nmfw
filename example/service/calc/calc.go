// Code generated using Nats Micro Service Framework version 10333f4459fab1d0f306e802b6a126a035fd2fed

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/choria-io/fisk"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	model "github.com/ripienaar/nmfw/example/service"
	impl "github.com/ripienaar/nmfw/example/impl"
	"github.com/sirupsen/logrus"
)

var (
	contextName string
	promPort    uint
	debug       bool
	maxRecon    int
)

func main() {
	app := fisk.New("calc", "Micro Service powered by NATS Micro")
	app.Version("0.0.2")

	app.Flag("debug", "Log at debug level").Envar("DEBUG").UnNegatableBoolVar(&debug)

	run := app.Command("run", "Runs the service").Action(runAction)
	run.Flag("context", "NATS Context to use for connection").Envar("CONTEXT").Default("MICRO").StringVar(&contextName)
	run.Flag("port", "Prometheus port for statistics").Envar("PORT").UintVar(&promPort)
	run.Flag("max-recon", "Maximum reconnection attempts").Envar("MAX_RECON").Default("60").IntVar(&maxRecon)

	app.MustParseWithUsage(os.Args[1:])
}

func startPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", promPort), nil)
}

func interruptWatcher(ctx context.Context, cancel context.CancelFunc, log *logrus.Entry) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigs:
			time.AfterFunc(5*time.Second, func() {
				log.Error("Forcing shutdown after 5 second since interrupt")
				os.Exit(1)
			})

			log.Infof("Shutting down on %s", sig)
			cancel()

		case <-ctx.Done():
			return
		}
	}
}

func runAction(_ *fisk.ParseContext) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logrus.SetLevel(logrus.InfoLevel)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{DisableHTMLEscape: true})
	log := logrus.NewEntry(logrus.StandardLogger()).WithFields(logrus.Fields{
		"service": "calc",
		"version": "0.0.2",
	})

	log.Infof("Starting calc service based on nmfw version 10333f4459fab1d0f306e802b6a126a035fd2fed")

	go interruptWatcher(ctx, cancel, log)

	if promPort > 0 {
		log.Infof("Starting Prometheus listener on :%d/metrics", promPort)
		go startPrometheus()
	}

	nc, err := natscontext.Connect(contextName,
		nats.MaxReconnects(maxRecon),
		nats.ConnectHandler(func(conn *nats.Conn) {
			log.Infof("Connected to %s", conn.ConnectedUrlRedacted())
		}),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			log.Infof("Reconnected to %s", conn.ConnectedUrlRedacted())
		}),
		nats.ErrorHandler(func(conn *nats.Conn, subscription *nats.Subscription, err error) {
			if err == nil {
				err = conn.LastError()
			}
			log.Errorf("NATS error encountered: %v", err)
		}),
		nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
			if err == nil {
				err = conn.LastError()
			}
			log.Errorf("NATS disconnected: %v", err)
		}),
		nats.ClosedHandler(func(conn *nats.Conn) {
			log.Errorf("NATS connection closed: %v", conn.LastError())
			cancel()
		}),
	)
	if err != nil {
		return err
	}

	svc := &model.CalcService{
		Name:        "calc",
		Version:     "0.0.2",
		Description: "Sample calculator service",
		RootSubject: "nmfw.calc",
		Log:         log,

		AverageHandler:    impl.AverageHandler,
		AddHandler:        impl.AddHandler,
		ExpressionHandler: impl.ExpressionHandler,
	}

	return svc.Start(ctx, nc)
}
