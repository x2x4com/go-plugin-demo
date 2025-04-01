package dynamic_plugin_shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type DynamicFunc struct {
	Name       string
	Call       func(args, options []interface{}) (interface{}, error)
	Help       string
	HasArgs    bool
	HasOptions bool
}

// func (f *DynamicFunc) SafeCall(args, options []interface{}) ([]interface{}, error) {
// 	val := reflect.ValueOf(f.Call)
// 	if val.Kind() != reflect.Func {
// 		panic("Call 字段必须为函数类型")
// 	}
//
// 	// 准备参数（需要实际参数类型严格匹配）
// 	in := make([]reflect.Value, len(args))
// 	for i, arg := range args {
// 		in[i] = reflect.ValueOf(arg)
// 	}
//
// 	// 调用函数并返回结果
// 	results := val.Call(in)
// 	out := make([]interface{}, len(results))
// 	for i, v := range results {
// 		out[i] = v.Interface()
// 	}
// 	return out
// }

func (f *DynamicFunc) GetFuncHelp() string {
	return f.Help
}

type DynamicPluginRPCClient struct {
	client *rpc.Client
}

func (c *DynamicPluginRPCClient) Invoke(method string, args, options []interface{}) (interface{}, error) {
	var resp interface{}
	err := c.client.Call("Plugin.Invoke", struct {
		Method  string
		Args    []interface{}
		Options []interface{}
	}{method, args, options}, &resp)
	return resp, err
}

func (c *DynamicPluginRPCClient) Help(method string) (string, error) {
	var resp string
	err := c.client.Call("Plugin.Help", method, &resp)
	return resp, err
}

func (c *DynamicPluginRPCClient) Version() string {
	var version string
	err := c.client.Call("Plugin.Version", struct{}{}, &version)
	if err != nil {
		return "unknown"
	}
	return version
}

type DynamicPluginRPCServer struct {
	Impl DynamicPluginInterface
}

func (s *DynamicPluginRPCServer) Invoke(args struct {
	Method  string
	Args    []interface{}
	Options []interface{}
}, resp *interface{}) error {
	result, err := s.Impl.Invoke(args.Method, args.Args, args.Options)
	*resp = result
	return err
}

func (s *DynamicPluginRPCServer) Help(method string, resp *string) error {
	result, err := s.Impl.Help(method)
	*resp = result
	return err
}

func (s *DynamicPluginRPCServer) Version(args interface{}, resp *string) error {
	*resp = s.Impl.Version()
	return nil
}

type DynamicPlugin struct {
	Impl DynamicPluginInterface
}

func (p *DynamicPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &DynamicPluginRPCServer{Impl: p.Impl}, nil
}

func (DynamicPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &DynamicPluginRPCClient{client: c}, nil
}
