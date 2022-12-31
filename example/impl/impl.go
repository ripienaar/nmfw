package impl

import (
	"fmt"
	"math"

	"github.com/antonmedv/expr"
	"github.com/ripienaar/nmfw/example/service"
)

func AverageHandler(req service.AverageRequest) (*service.CalcResponse, error) {
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

func AddHandler(req service.AddRequest) (*service.CalcResponse, error) {
	resp := service.CalcResponse{Operation: "add"}
	if len(req.Values) == 0 {
		return &resp, nil
	}

	for _, v := range req.Values {
		resp.Result += v
	}

	return &resp, nil
}

func ExpressionHandler(req service.ExpressionRequest) (*service.CalcResponse, error) {
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
