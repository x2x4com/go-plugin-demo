package main

import (
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
	handshake := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "DYNAMIC_PLUGIN_STRING_UTILS",
		MagicCookieValue: "string_utils",
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshake,
		Plugins: map[string]plugin.Plugin{
			"string_utils": &StringUtilsPlugin{
				Impl: &StringUtilsImplementation{},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
