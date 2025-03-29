package main

import (
	"go-plugin-demo/src/shared"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// StringUtilsImplementation 实现字符串工具接口
type StringUtilsImplementation struct{}

func (s *StringUtilsImplementation) Reverse(str string) (string, error) {
	return Reverse(str), nil
}

func (s *StringUtilsImplementation) ToUpper(str string) (string, error) {
	return ToUpper(str), nil
}

func (s *StringUtilsImplementation) ToLower(str string) (string, error) {
	return ToLower(str), nil
}

func (s *StringUtilsImplementation) ToTitle(str string) (string, error) {
	return ToTitle(str), nil
}

func (s *StringUtilsImplementation) ToCamel(str string) (string, error) {
	return ToCamel(str), nil
}

func (s *StringUtilsImplementation) ToSnake(str string) (string, error) {
	return ToSnake(str), nil
}

// StringUtilsPlugin 实现plugin.Plugin接口
type StringUtilsPlugin struct {
	Impl *StringUtilsImplementation
}

func (p *StringUtilsPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return p.Impl, nil
}

func (p *StringUtilsPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return p.Impl, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"string_utils": &StringUtilsPlugin{
				Impl: &StringUtilsImplementation{},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
