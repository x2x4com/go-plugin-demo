package cmd

import (
	"encoding/json"
	"fmt"
	"go-plugin-demo/src/shared"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

type PluginConfig struct {
	Name      string                 `json:"name"`
	Path      string                 `json:"path"`
	Handshake plugin.HandshakeConfig `json:"handshake"`
}

type Config struct {
	Plugins []PluginConfig `json:"plugins"`
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有插件及方法",
	Long:  "列出当前系统中所有已注册的插件及其可用方法",
	Run: func(cmd *cobra.Command, args []string) {
		// 读取插件配置
		configFile, err := os.ReadFile("config/plugins.json")
		if err != nil {
			color.Red("读取插件配置文件失败: %v", err)
			return
		}

		var config Config
		if err := json.Unmarshal(configFile, &config); err != nil {
			color.Red("解析插件配置失败: %v", err)
			return
		}

		// 初始化彩色输出
		blue := color.New(color.FgBlue).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		fmt.Println(blue("\n已注册插件:"))
		for _, plugin := range config.Plugins {
			fmt.Printf("\n%s %s\n", green("插件名称:"), plugin.Name)
			fmt.Printf("%s %s\n", green("插件路径:"), plugin.Path)

			// 获取插件ABI信息
			abi := getPluginABI(plugin.Name)
			if abi != nil {
				fmt.Println(yellow("\n可用方法:"))
				for method, spec := range abi.Methods {
					fmt.Printf("  %s(", method)
					for i, param := range spec.Params {
						if i > 0 {
							fmt.Print(", ")
						}
						fmt.Print(param)
					}
					fmt.Printf(") → %s\n", spec.Returns)
				}
			}
		}
	},
}

func getPluginABI(name string) *shared.PluginABI {
	// 这里简化处理，实际应从插件获取ABI信息
	switch name {
	case "calculator":
		return shared.CalculatorABI()
	case "string_utils":
		return &shared.PluginABI{
			Name:    "string_utils",
			Version: "1.0.0",
			Methods: map[string]shared.MethodSpec{
				"Reverse": {Params: []string{"string"}, Returns: "string"},
				"ToUpper": {Params: []string{"string"}, Returns: "string"},
				"ToCamel": {Params: []string{"string"}, Returns: "string"},
			},
		}
	case "date_utils":
		return &shared.PluginABI{
			Name:    "date_utils",
			Version: "1.0.0",
			Methods: map[string]shared.MethodSpec{
				"AddDays": {Params: []string{"time.Time", "int"}, Returns: "time.Time"},
				"Format":  {Params: []string{"time.Time", "string"}, Returns: "string"},
				"Between": {Params: []string{"time.Time", "time.Time"}, Returns: "int"},
			},
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
