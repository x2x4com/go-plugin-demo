package main

import (
	"net/rpc"
	"time"

	"github.com/hashicorp/go-plugin"
)

// DateUtilsImplementation 实现日期工具接口
type DateUtilsImplementation struct{}

func (d *DateUtilsImplementation) AddDays(date time.Time, days int) (time.Time, error) {
	return AddDays(date, days), nil
}

func (d *DateUtilsImplementation) Format(date time.Time, layout string) (string, error) {
	return Format(date, layout), nil
}

func (d *DateUtilsImplementation) Parse(dateStr, layout string) (time.Time, error) {
	return Parse(dateStr, layout)
}

func (d *DateUtilsImplementation) Between(start, end time.Time) (int, error) {
	return Between(start, end), nil
}

// DateUtilsPlugin 实现plugin.Plugin接口
type DateUtilsPlugin struct {
	Impl *DateUtilsImplementation
}

func (p *DateUtilsPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return p.Impl, nil
}

func (p *DateUtilsPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return p.Impl, nil
}

func main() {
	handshake := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "DYNAMIC_PLUGIN_DATE_UTILS",
		MagicCookieValue: "date_utils",
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshake,
		Plugins: map[string]plugin.Plugin{
			"date_utils": &DateUtilsPlugin{
				Impl: &DateUtilsImplementation{},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
