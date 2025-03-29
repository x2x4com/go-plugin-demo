package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go-plugin-demo/src/shared"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-plugin"
)

// 插件配置结构
type PluginConfig struct {
	Name      string                 `json:"name"`
	Path      string                 `json:"path"`
	Handshake plugin.HandshakeConfig `json:"handshake"`
}

// 配置文件结构
type Config struct {
	Plugins []PluginConfig `json:"plugins"`
}

func main() {
	// 1. 读取配置文件
	configFile, err := os.ReadFile("config/plugins.json")
	if err != nil {
		log.Fatal("读取插件配置文件失败:", err)
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Fatal("解析插件配置失败:", err)
	}

	// 2. 初始化插件管理器
	pm := NewPluginManager()
	defer pm.UnloadAll()

	// 3. 加载所有插件
	for _, pluginConfig := range config.Plugins {
		// 校验插件路径
		if _, err := os.Stat(pluginConfig.Path); os.IsNotExist(err) {
			log.Printf("插件 %s 路径不存在: %s", pluginConfig.Name, pluginConfig.Path)
			continue
		}

		// 转换为绝对路径
		absPath, err := filepath.Abs(pluginConfig.Path)
		if err != nil {
			log.Printf("获取插件 %s 绝对路径失败: %v", pluginConfig.Name, err)
			continue
		}
		pluginConfig.Path = absPath

		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: shared.Handshake,
			Plugins: map[string]plugin.Plugin{
				pluginConfig.Name: &shared.DynamicPluginRPC{},
			},
			Cmd: exec.Command(pluginConfig.Path),
		})

		// 连接RPC客户端
		rpcClient, err := client.Client()
		if err != nil {
			log.Printf("连接插件 %s RPC失败: %v", pluginConfig.Name, err)
			continue
		}

		// 获取插件实例
		if _, err := rpcClient.Dispense(pluginConfig.Name); err != nil {
			log.Printf("获取插件 %s 实例失败: %v", pluginConfig.Name, err)
			continue
		}

		// 注册插件
		pm.plugins[pluginConfig.Name] = client
		fmt.Printf("成功加载插件: %s\n", pluginConfig.Name)
	}

	// 4. 交互式菜单
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n请选择要测试的插件:")
		fmt.Println("1. 计算器插件")
		fmt.Println("2. 字符串插件")
		fmt.Println("3. 日期插件")
		fmt.Println("4. 退出")
		fmt.Print("请输入选项: ")

		input, _ := reader.ReadString('\n')
		switch input {
		case "1\n":
			if _, ok := pm.plugins["calculator"]; ok {
				testCalculator(pm)
			} else {
				fmt.Println("计算器插件未加载")
			}
		case "2\n":
			if _, ok := pm.plugins["string_utils"]; ok {
				testStringUtils(pm)
			} else {
				fmt.Println("字符串插件未加载")
			}
		case "3\n":
			if _, ok := pm.plugins["date_utils"]; ok {
				testDateUtils(pm)
			} else {
				fmt.Println("日期插件未加载")
			}
		case "4\n":
			return
		default:
			fmt.Println("无效选项")
		}
	}
}

func testCalculator(pm *PluginManager) {
	fmt.Println("\n测试计算器插件:")
	if sum, err := pm.Invoke("calculator", "Add", 5, 3); err == nil {
		fmt.Println("5 + 3 =", sum)
	}
	if diff, err := pm.Invoke("calculator", "Subtract", 5, 3); err == nil {
		fmt.Println("5 - 3 =", diff)
	}
	if product, err := pm.Invoke("calculator", "Multiply", 5, 3); err == nil {
		fmt.Println("5 * 3 =", product)
	}
	if quotient, err := pm.Invoke("calculator", "Divide", 5, 3); err == nil {
		fmt.Printf("5 / 3 = %.2f\n", quotient)
	}
}

func testStringUtils(pm *PluginManager) {
	fmt.Println("\n测试字符串插件:")
	if result, err := pm.Invoke("string_utils", "Reverse", "hello"); err == nil {
		fmt.Println("Reverse('hello') =", result)
	}
	if result, err := pm.Invoke("string_utils", "ToUpper", "hello"); err == nil {
		fmt.Println("ToUpper('hello') =", result)
	}
	if result, err := pm.Invoke("string_utils", "ToCamel", "hello world"); err == nil {
		fmt.Println("ToCamel('hello world') =", result)
	}
}

func testDateUtils(pm *PluginManager) {
	fmt.Println("\n测试日期插件:")
	now := time.Now()
	if result, err := pm.Invoke("date_utils", "AddDays", now, 7); err == nil {
		fmt.Printf("AddDays(%v, 7) = %v\n", now.Format("2006-01-02"), result.(time.Time).Format("2006-01-02"))
	}
	if result, err := pm.Invoke("date_utils", "Format", now, "2006-01-02 15:04:05"); err == nil {
		fmt.Println("Format(now) =", result)
	}
	if result, err := pm.Invoke("date_utils", "Between", now, now.AddDate(0, 0, 14)); err == nil {
		fmt.Println("Between(now, now+14d) =", result, "days")
	}
}
