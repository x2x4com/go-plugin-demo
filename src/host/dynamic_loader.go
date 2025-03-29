package main

import (
	"fmt"
	"go-plugin-demo/src/shared"
	"log"
	"os"
	"os/exec"

	goplugin "github.com/hashicorp/go-plugin"
)

// PluginManager 管理动态加载的插件
type PluginManager struct {
	plugins map[string]*goplugin.Client
	abis    map[string]*shared.PluginABI
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]*goplugin.Client),
		abis:    make(map[string]*shared.PluginABI),
	}
}

// LoadPlugins 从配置目录加载所有插件
func (pm *PluginManager) LoadPlugins(cfg *PluginConfig) error {
	info, err := os.Stat(cfg.Path)
	if info.IsDir() {
		return fmt.Errorf("%s必须是一个文件而不是目录", cfg.Path)
	} else if err != nil {
		return fmt.Errorf("读取插件失败: %v", err)
	}

	handshake := goplugin.HandshakeConfig{
		ProtocolVersion:  cfg.Handshake.ProtocolVersion,
		MagicCookieKey:   cfg.Handshake.MagicCookieKey,
		MagicCookieValue: cfg.Handshake.MagicCookieValue,
	}

	// 1. 创建插件客户端
	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: handshake,
		Plugins:         map[string]goplugin.Plugin{cfg.Name: &shared.DynamicPluginRPC{}},
		Cmd:             exec.Command(cfg.Path),
	})

	// 2. 连接RPC客户端
	rpcClient, err := client.Client()
	if err != nil {
		return fmt.Errorf("RPC连接失败: %v", err)
	}

	// 3. 获取插件实例
	raw, err := rpcClient.Dispense(cfg.Name)
	if err != nil {
		return fmt.Errorf("获取插件实例失败: %v", err)
	}

	fmt.Printf("插件实例: %v\n", raw)

	// 4. 获取ABI描述
	//dynamicPlugin := raw.(shared.DynamicPlugin)
	//abi, err := dynamicPlugin.GetABI()
	//if err != nil {
	//	return fmt.Errorf("获取ABI失败: %v", err)
	//}

	pm.plugins[cfg.Name] = client
	//pm.abis[abi.Name] = abi
	log.Printf("成功加载插件: %s v%s", cfg.Name, "aa")

	return nil
}

// Invoke 动态调用插件方法
func (pm *PluginManager) Invoke(pluginName, method string, args ...interface{}) (interface{}, error) {
	client, ok := pm.plugins[pluginName]
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
	plugin := raw.(shared.DynamicPlugin)
	return plugin.Invoke(method, args...)
}

// UnloadAll 卸载所有插件
func (pm *PluginManager) UnloadAll() {
	for _, client := range pm.plugins {
		client.Kill()
	}
}
