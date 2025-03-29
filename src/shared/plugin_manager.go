package shared

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-plugin"
	goplugin "github.com/hashicorp/go-plugin"
)

// Handshake 插件握手配置
var Handshake = goplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "DYNAMIC_PLUGIN",
	MagicCookieValue: "dynamic",
}

// PluginManager 管理动态加载的插件
type PluginManager struct {
	Plugins map[string]*goplugin.Client
	ABIs    map[string]*PluginABI
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		Plugins: make(map[string]*goplugin.Client),
		ABIs:    make(map[string]*PluginABI),
	}
}

// LoadFromConfig 从配置文件加载插件
func (pm *PluginManager) LoadFromConfig(configPath string) error {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取插件配置文件失败: %v", err)
	}

	var config struct {
		Plugins []struct {
			Name      string                 `json:"name"`
			Path      string                 `json:"path"`
			Handshake plugin.HandshakeConfig `json:"handshake"`
		} `json:"plugins"`
	}

	if err := json.Unmarshal(configFile, &config); err != nil {
		return fmt.Errorf("解析插件配置失败: %v", err)
	}

	for _, pluginConfig := range config.Plugins {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: pluginConfig.Handshake,
			Plugins:         map[string]plugin.Plugin{"dynamic": &DynamicPluginRPC{}},
			Cmd:             exec.Command(pluginConfig.Path),
		})

		rpcClient, err := client.Client()
		if err != nil {
			log.Printf("连接插件 %s RPC失败: %v", pluginConfig.Name, err)
			continue
		}

		raw, err := rpcClient.Dispense("dynamic")
		if err != nil {
			log.Printf("获取插件 %s 实例失败: %v", pluginConfig.Name, err)
			continue
		}

		plugin := raw.(DynamicPlugin)
		generator := NewABIGenerator()
		abi, err := generator.GenerateFromInstance(pluginConfig.Name, "1.0.0", plugin)
		if err != nil {
			log.Printf("生成插件 %s ABI失败: %v", pluginConfig.Name, err)
			continue
		}

		pm.Plugins[pluginConfig.Name] = client
		pm.ABIs[pluginConfig.Name] = abi
		log.Printf("成功加载插件: %s", pluginConfig.Name)
	}

	return nil
}

func (pm *PluginManager) loadPlugin(path string) (*goplugin.Client, *PluginABI, error) {
	// 1. 创建插件客户端
	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         map[string]goplugin.Plugin{"dynamic": &DynamicPluginRPC{}},
		Cmd:             exec.Command(path),
	})

	// 2. 连接RPC客户端
	rpcClient, err := client.Client()
	if err != nil {
		return nil, nil, fmt.Errorf("RPC连接失败: %v", err)
	}

	// 3. 获取插件实例
	raw, err := rpcClient.Dispense("dynamic")
	if err != nil {
		return nil, nil, fmt.Errorf("获取插件实例失败: %v", err)
	}

	// 4. 生成ABI描述
	dynamicPlugin := raw.(DynamicPlugin)

	// 检查缓存
	if cachedABI, ok := pm.ABIs[filepath.Base(path)]; ok {
		return client, cachedABI, nil
	}

	// 使用ABI生成器
	generator := NewABIGenerator()
	abi, err := generator.GenerateFromInstance(
		filepath.Base(path),
		"1.0.0",
		dynamicPlugin,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("生成ABI失败: %v", err)
	}

	// 存入缓存
	pm.ABIs[filepath.Base(path)] = abi

	return client, abi, nil
}

// Invoke 动态调用插件方法
func (pm *PluginManager) Invoke(pluginName, method string, args ...interface{}) (interface{}, error) {
	client, ok := pm.Plugins[pluginName]
	if !ok {
		return nil, fmt.Errorf("插件 %s 未加载", pluginName)
	}

	// 获取RPC客户端
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	// 动态调用
	raw, err := rpcClient.Dispense("dynamic")
	if err != nil {
		return nil, err
	}

	// 类型断言并调用
	plugin := raw.(DynamicPlugin)
	return plugin.Invoke(method, args...)
}

// UnloadAll 卸载所有插件
func (pm *PluginManager) UnloadAll() {
	for _, client := range pm.Plugins {
		client.Kill()
	}
}
