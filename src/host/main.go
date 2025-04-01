package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	dynamic_plugin_shared "go-plugin-demo/src/internal/plugin/shared"
	calculator "go-plugin-demo/src/plugins/calculator/shared"
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
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	// 1. 读取配置文件
	configFile, err := os.ReadFile("config/plugins.json")
	if err != nil {
		logger.Error("读取插件配置文件失败:", err)
		return
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		logger.Error("解析插件配置失败:", err)
		return
	}

	// // 2. 初始化插件管理器
	pm := NewPluginManager()
	defer pm.UnloadAll()
	//
	// // 3. 加载所有插件
	// for _, pluginConfig := range config.Plugins {
	// 	err := pm.LoadPlugins(&pluginConfig)
	// 	if err != nil {
	// 		log.Fatal(pluginConfig.Name, " 加载插件失败: ", err)
	// 		continue
	// 	}
	// }

	handshake := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "DYNAMIC_PLUGIN_CALCULATOR",
		MagicCookieValue: "calculator",
	}

	// pluginMap is the map of plugins we can dispense.
	pluginMap := map[string]plugin.Plugin{
		"calculator": &calculator.CalculatorPlugin{},
	}

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshake,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./bin/plugins/calculator"),
		Logger:          logger,
	})

	rpcClient, err := client.Client()
	if err != nil {
		logger.Error(fmt.Sprintf("RPC连接失败: %v", err))
		return
	}

	// 3. 获取插件实例
	raw, err := rpcClient.Dispense("calculator")
	if err != nil {
		logger.Error(fmt.Sprintf("获取插件实例失败: %v", err))
	}

	fmt.Printf("插件实例: %v\n", raw)

	pm.plugins["calculator"] = client

	handshake = dynamic_plugin_shared.GenHandShakeConfig("calculator")

	pluginMap = map[string]plugin.Plugin{
		"date_utils": &dynamic_plugin_shared.DynamicPlugin{},
	}

	client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshake,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./bin/plugins/date_utils"),
		Logger:          logger,
	})

	rpcClient, err = client.Client()
	if err != nil {
		logger.Error(fmt.Sprintf("RPC连接失败: %v", err))
		return
	}

	raw, err = rpcClient.Dispense("date_utils")
	if err != nil {
		logger.Error(fmt.Sprintf("获取插件实例失败: %v", err))
	}

	fmt.Printf("插件实例: %v\n", raw)

	pm.plugins["date_utils"] = client

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
				rpcClient, _ := pm.plugins["calculator"].Client()
				// 3. 获取插件实例
				raw, err := rpcClient.Dispense("calculator")
				if err != nil {
					logger.Error(fmt.Sprintf("获取插件实例失败: %v", err))
				}
				calc := raw.(calculator.Calculator)
				res, err := calc.Add(1, 2)
				fmt.Printf("计算结果: %f\nerr: %v", res, err)
				res, err = calc.Subtract(1, 2)
				fmt.Printf("计算结果: %f\nerr: %v", res, err)
				res, err = calc.Multiply(1, 2)
				fmt.Printf("计算结果: %f\nerr: %v", res, err)
				res, err = calc.Divide(1, 2)
				fmt.Printf("计算结果: %f\nerr: %v", res, err)

			} else {
				fmt.Println("计算器插件未加载")
			}
		//case "2\n":
		//	if _, ok := pm.plugins["string_utils"]; ok {
		//		testStringUtils(pm)
		//	} else {
		//		fmt.Println("字符串插件未加载")
		//	}
		case "3\n":
			if _, ok := pm.plugins["date_utils"]; ok {
				rpcClient, _ := pm.plugins["date_utils"].Client()
				// 3. 获取插件实例
				raw, err := rpcClient.Dispense("date_utils")
				if err != nil {
					logger.Error(fmt.Sprintf("获取插件实例失败: %v", err))
				}
				dp := raw.(dynamic_plugin_shared.DynamicPluginInterface)
				n := time.Now().Format(time.RFC3339)
				fmt.Println("当前时间:", n)
				res, err := dp.Invoke("AddDays", []interface{}{n, 5}, []interface{}{})
				fmt.Printf("计算结果: %v\nerr: %v", res, err)
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

/*
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
*/
