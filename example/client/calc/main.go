package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/choria-io/fisk"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/ripienaar/nmfw/example/service"
)

var (
	contextName string
	values      []float32
	expression  string
	timeOut     time.Duration
)

func main() {
	app := fisk.New("calc", "NATS Calculator Service Client")
	app.Flag("context", "NATS Context to use for connection").Envar("CONTEXT").Default("MICRO").StringVar(&contextName)
	app.Flag("timeout", "How long to wait for responses").Envar("TIMEOUT").Default("5s").DurationVar(&timeOut)

	add := app.Command("add", "Adds a series of numbers").Alias("+").Action(addAction)
	add.Arg("numbers", "Numbers to add").Float32ListVar(&values)

	avg := app.Command("average", "Calculates the average of a series of numbers").Alias("avg").Action(avgAction)
	avg.Arg("numbers", "Numbers to average").Float32ListVar(&values)

	expr := app.Command("expr", "Calculates an expression").Alias("e").Action(expAction)
	expr.Arg("expression", "The expression to calculate").StringVar(&expression)

	app.MustParseWithUsage(os.Args[1:])
}

func client(cb func(client *service.CalcClient) error) error {
	nc, err := natscontext.Connect(contextName)
	if err != nil {
		return err
	}

	c := service.NewCalcClient(nc, timeOut)

	return cb(c)
}

func expAction(_ *fisk.ParseContext) error {
	return client(func(c *service.CalcClient) error {
		res, err := c.Expression(context.Background(), service.ExpressionRequest{Expression: expression})
		if err != nil {
			return err
		}

		fmt.Println(res.Result)

		return nil
	})
}

func addAction(_ *fisk.ParseContext) error {
	return client(func(c *service.CalcClient) error {

		res, err := c.Add(context.Background(), service.AddRequest{Values: values})
		if err != nil {
			return err
		}

		fmt.Println(res.Result)

		return nil
	})
}

func avgAction(_ *fisk.ParseContext) error {
	return client(func(c *service.CalcClient) error {
		res, err := c.Average(context.Background(), service.AverageRequest{Values: values})
		if err != nil {
			return err
		}

		fmt.Println(res.Result)

		return nil
	})
}
