package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
		err := pm.LoadPlugins(&pluginConfig)
		if err != nil {
			log.Fatal(pluginConfig.Name, " 加载插件失败: ", err)
			continue
		}
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
