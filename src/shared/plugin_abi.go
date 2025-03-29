package shared

import (
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// PluginABI 描述插件的接口规范
type PluginABI struct {
	Name    string                `json:"name"`
	Version string                `json:"version"`
	Methods map[string]MethodSpec `json:"methods"`
}

// MethodSpec 描述方法签名
type MethodSpec struct {
	Params  []string `json:"params"`
	Returns string   `json:"returns"`
}

// PluginDescriptor 插件描述文件结构
type PluginDescriptor struct {
	Path string    `json:"path"`
	ABI  PluginABI `json:"abi"`
}

// DynamicPlugin 动态插件接口
type DynamicPlugin interface {
	plugin.Plugin
	GetABI() (*PluginABI, error)
	Invoke(method string, args ...interface{}) (interface{}, error)
}

// DynamicPluginRPC RPC实现
type DynamicPluginRPC struct {
	Impl DynamicPlugin
}

func (p *DynamicPluginRPC) Server(*plugin.MuxBroker) (interface{}, error) {
	return p.Impl, nil
}

func (p *DynamicPluginRPC) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &DynamicPluginRPC{Impl: p.Impl}, nil
}

func (p *DynamicPluginRPC) GetABI(args interface{}, resp *PluginABI) error {
	// 自动生成ABI
	generator := NewABIGenerator()
	abi, err := generator.GenerateFromInstance("plugin", "1.0.0", p.Impl)
	if err != nil {
		return err
	}
	*resp = *abi
	return nil
}

func (p *DynamicPluginRPC) Invoke(args []interface{}, resp *interface{}) error {
	if len(args) < 1 {
		return fmt.Errorf("至少需要方法名参数")
	}
	method, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("第一个参数必须是方法名")
	}

	result, err := p.Impl.Invoke(method, args[1:]...)
	if err != nil {
		return err
	}
	*resp = result
	return nil
}

// CalculatorABI 生成计算器插件的ABI描述
func CalculatorABI() *PluginABI {
	return &PluginABI{
		Name:    "calculator",
		Version: "1.0.0",
		Methods: map[string]MethodSpec{
			"Add":      {Params: []string{"int", "int"}, Returns: "int,error"},
			"Subtract": {Params: []string{"int", "int"}, Returns: "int,error"},
			"Multiply": {Params: []string{"int", "int"}, Returns: "int,error"},
			"Divide":   {Params: []string{"int", "int"}, Returns: "float64,error"},
		},
	}
}
