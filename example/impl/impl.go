package impl

import (
	"context"
	"fmt"
	"math"

	"github.com/antonmedv/expr"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/ripienaar/nmfw/example/service"
	"github.com/sirupsen/logrus"
)

type request interface {
	Logger() *logrus.Entry
	Conn() *nats.Conn
	Request() *micro.Request
}

func AverageHandler(_ context.Context, req service.AverageRequest) (*service.CalcResponse, error) {
	resp := service.CalcResponse{Operation: "average"}
	if len(req.Values) == 0 {
		return &resp, nil
	}

	for _, v := range req.Values {
		resp.Result += v
	}

	resp.Result = resp.Result / float32(len(req.Values))

	return &resp, nil
}

func AddHandler(_ context.Context, req service.AddRequest) (*service.CalcResponse, error) {
	resp := service.CalcResponse{Operation: "add"}
	if len(req.Values) == 0 {
		return &resp, nil
	}

	for _, v := range req.Values {
		resp.Result += v
	}

	return &resp, nil
}

func ExpressionHandler(ctx context.Context, req service.ExpressionRequest) (*service.CalcResponse, error) {
	helper := ctx.Value("nmfw").(request)
	helper.Logger().Infof("Calculating expression %s", req.Expression)

	program, err := expr.Compile(req.Expression)
	if err != nil {
		return nil, err
	}

	output, err := expr.Run(program, nil)
	if err != nil {
		return nil, err
	}

	res := &service.CalcResponse{Operation: "expr"}

	switch n := output.(type) {
	case int:
		res.Result = float32(n)
	case float64:
		if n > math.MaxFloat32 {
			return nil, fmt.Errorf("result too big")
		}

		res.Result = float32(n)
	default:
		return nil, fmt.Errorf("unsupported expression result type")
	}

	return res, nil
}
