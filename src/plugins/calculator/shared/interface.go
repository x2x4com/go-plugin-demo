package shared

import (
	"errors"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Calculator interface {
	Add(a, b float64) (float64, error)
	Subtract(a, b float64) (float64, error)
	Multiply(a, b float64) (float64, error)
	Divide(a, b float64) (float64, error)
}

type CalculatorRPCClient struct{ client *rpc.Client }

func (c *CalculatorRPCClient) Add(a, b float64) (float64, error) {
	var resp float64
	args := struct{ A, B float64 }{a, b}
	err := c.client.Call("Plugin.Add", args, &resp)
	return resp, err
}

func (c *CalculatorRPCClient) Subtract(a, b float64) (float64, error) {
	var resp float64
	args := struct{ A, B float64 }{a, b}
	err := c.client.Call("Plugin.Subtract", args, &resp)
	return resp, err
}

func (c *CalculatorRPCClient) Multiply(a, b float64) (float64, error) {
	var resp float64
	args := struct{ A, B float64 }{a, b}
	err := c.client.Call("Plugin.Multiply", args, &resp)
	return resp, err
}

func (c *CalculatorRPCClient) Divide(a, b float64) (float64, error) {
	var resp float64
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	args := struct{ A, B float64 }{a, b}
	err := c.client.Call("Plugin.Divide", args, &resp)
	return resp, err
}

type CalculatorRPCServer struct {
	Impl Calculator
}

func (s *CalculatorRPCServer) Add(args *struct{ A, B float64 }, resp *float64) error {
	v, err := s.Impl.Add(args.A, args.B)
	*resp = v
	return err
}

func (s *CalculatorRPCServer) Subtract(args *struct{ A, B float64 }, resp *float64) error {
	v, err := s.Impl.Subtract(args.A, args.B)
	*resp = v
	return err
}

func (s *CalculatorRPCServer) Multiply(args *struct{ A, B float64 }, resp *float64) error {
	v, err := s.Impl.Multiply(args.A, args.B)
	*resp = v
	return err
}

func (s *CalculatorRPCServer) Divide(args *struct{ A, B float64 }, resp *float64) error {
	v, err := s.Impl.Divide(args.A, args.B)
	*resp = v
	return err
}

type CalculatorPlugin struct {
	Impl Calculator
}

func (p *CalculatorPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &CalculatorRPCServer{Impl: p.Impl}, nil
}

func (CalculatorPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CalculatorRPCClient{client: c}, nil
}
