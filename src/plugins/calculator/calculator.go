package main

import (
	"errors"
	"go-plugin-demo/src/plugins/calculator/shared"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type CalculatorImpl struct {
	logger hclog.Logger
}

func (c *CalculatorImpl) Add(a, b float64) (float64, error) {
	res := a + b
	c.logger.Info("Add called", "a", a, "b", b, "result", res)
	return res, nil
}

func (c *CalculatorImpl) Subtract(a, b float64) (float64, error) {
	res := a - b
	c.logger.Info("Subtract called", "a", a, "b", b, "result", res)
	return res, nil
}

func (c *CalculatorImpl) Multiply(a, b float64) (float64, error) {
	res := a * b
	c.logger.Info("Multiply called", "a", a, "b", b, "result", res)
	return res, nil
}

func (c *CalculatorImpl) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	res := a / b
	c.logger.Info("Divide called", "a", a, "b", b, "result", res)
	return res, nil
}

var handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "DYNAMIC_PLUGIN_CALCULATOR",
	MagicCookieValue: "calculator",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	calc := &CalculatorImpl{logger: logger}

	var pluginMap = map[string]plugin.Plugin{
		"calculator": &shared.CalculatorPlugin{Impl: calc},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshake,
		Plugins:         pluginMap,
	})

}
