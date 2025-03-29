package shared

import (
	"fmt"
	"reflect"
	"sync"
)

// ABIGenerator 自动生成插件ABI描述
type ABIGenerator struct {
	typeMap map[string]string
	mu      sync.RWMutex
}

// 基础类型映射表
var basicTypeMap = map[string]string{
	"int":       "int64",
	"int8":      "int32",
	"int16":     "int32",
	"int32":     "int32",
	"int64":     "int64",
	"uint":      "uint64",
	"uint8":     "uint32",
	"uint16":    "uint32",
	"uint32":    "uint32",
	"uint64":    "uint64",
	"float32":   "float64",
	"float64":   "float64",
	"complex64": "complex128",
}

// NewABIGenerator 创建ABI生成器实例
func NewABIGenerator() *ABIGenerator {
	return &ABIGenerator{
		typeMap: make(map[string]string),
	}
}

// AddTypeMapping 添加自定义类型映射
func (g *ABIGenerator) AddTypeMapping(from, to string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.typeMap[from] = to
}

// GetTypeMapping 获取类型映射
func (g *ABIGenerator) GetTypeMapping(from string) (string, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	val, ok := g.typeMap[from]
	return val, ok
}

// GenerateFromInstance 从插件实例生成ABI描述
func (g *ABIGenerator) GenerateFromInstance(pluginName, version string, plugin interface{}) (*PluginABI, error) {
	abi := &PluginABI{
		Name:    pluginName,
		Version: version,
		Methods: make(map[string]MethodSpec),
	}

	pluginType := reflect.TypeOf(plugin)
	if pluginType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("插件实例必须是指针类型")
	}

	// 遍历所有导出的方法
	for i := 0; i < pluginType.NumMethod(); i++ {
		method := pluginType.Method(i)
		if !method.IsExported() {
			continue
		}

		methodSpec := MethodSpec{
			Params:  make([]string, method.Type.NumIn()-1), // 减去接收者参数
			Returns: g.getTypeName(method.Type.Out(0)),
		}

		// 处理参数类型
		for j := 1; j < method.Type.NumIn(); j++ {
			methodSpec.Params[j-1] = g.getTypeName(method.Type.In(j))
		}

		// 处理错误返回值
		if method.Type.NumOut() > 1 {
			methodSpec.Returns += ",error"
		}

		abi.Methods[method.Name] = methodSpec
	}

	return abi, nil
}

// getTypeName 获取类型的可读名称
func (g *ABIGenerator) getTypeName(t reflect.Type) string {
	// 检查基础类型映射
	if mapped, ok := basicTypeMap[t.Name()]; ok {
		return mapped
	}

	// 检查自定义类型映射
	if mapped, ok := g.GetTypeMapping(t.Name()); ok {
		return mapped
	}

	switch t.Kind() {
	case reflect.Ptr:
		return "*" + g.getTypeName(t.Elem())
	case reflect.Slice:
		return "[]" + g.getTypeName(t.Elem())
	case reflect.Map:
		return "map[" + g.getTypeName(t.Key()) + "]" + g.getTypeName(t.Elem())
	case reflect.Interface:
		if t.Name() == "" {
			return "interface{}"
		}
		return t.Name()
	default:
		return t.Name()
	}
}

// GenerateFromPackage 从包路径生成ABI描述(需要插件实现DynamicPlugin接口)
func (g *ABIGenerator) GenerateFromPackage(pluginName, version, pkgPath string) (*PluginABI, error) {
	// TODO: 实现从包路径加载并分析
	return nil, fmt.Errorf("未实现")
}
