package main

import (
	"errors"
	"net/rpc"

	"go-plugin-demo/src/shared"

	"github.com/hashicorp/go-plugin"
)

// CalculatorImplementation 实现计算器接口
type CalculatorImplementation struct{}

func (c *CalculatorImplementation) Add(a, b int) (int, error) {
	return a + b, nil
}

func (c *CalculatorImplementation) Subtract(a, b int) (int, error) {
	return a - b, nil
}

func (c *CalculatorImplementation) Multiply(a, b int) (int, error) {
	return a * b, nil
}

func (c *CalculatorImplementation) Divide(a, b int) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return float64(a) / float64(b), nil
}

// CalculatorPlugin 实现DynamicPlugin接口
type CalculatorPlugin struct {
	Impl *CalculatorImplementation
}

func (p *CalculatorPlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return p.Impl, nil
}

func (p *CalculatorPlugin) Client(broker *plugin.MuxBroker, client *rpc.Client) (interface{}, error) {
	return &CalculatorRPC{client: client}, nil
}

func (p *CalculatorPlugin) GetABI() (*shared.PluginABI, error) {
	return shared.CalculatorABI(), nil
}

func (p *CalculatorPlugin) Invoke(method string, args ...interface{}) (interface{}, error) {
	switch method {
	case "Add":
		return p.Impl.Add(args[0].(int), args[1].(int))
	case "Subtract":
		return p.Impl.Subtract(args[0].(int), args[1].(int))
	case "Multiply":
		return p.Impl.Multiply(args[0].(int), args[1].(int))
	case "Divide":
		return p.Impl.Divide(args[0].(int), args[1].(int))
	default:
		return nil, errors.New("unknown method")
	}
}

// CalculatorRPC RPC客户端实现
type CalculatorRPC struct {
	client *rpc.Client
}

func (c *CalculatorRPC) Add(a, b int) (int, error) {
	var resp int
	err := c.client.Call("Plugin.Add", []int{a, b}, &resp)
	return resp, err
}

func (c *CalculatorRPC) Subtract(a, b int) (int, error) {
	var resp int
	err := c.client.Call("Plugin.Subtract", []int{a, b}, &resp)
	return resp, err
}

func (c *CalculatorRPC) Multiply(a, b int) (int, error) {
	var resp int
	err := c.client.Call("Plugin.Multiply", []int{a, b}, &resp)
	return resp, err
}

func (c *CalculatorRPC) Divide(a, b int) (float64, error) {
	var resp float64
	err := c.client.Call("Plugin.Divide", []int{a, b}, &resp)
	return resp, err
}

func main() {
	handshake := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "DYNAMIC_PLUGIN_CALCULATOR",
		MagicCookieValue: "calculator",
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshake,
		Plugins: map[string]plugin.Plugin{
			"calculator": &CalculatorPlugin{
				Impl: &CalculatorImplementation{},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
