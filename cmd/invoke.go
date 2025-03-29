package cmd

import (
	"fmt"
	"go-plugin-demo/src/shared"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var invokeCmd = &cobra.Command{
	Use:   "invoke [插件名] [方法名] [参数...]",
	Short: "调用插件方法",
	Long:  "调用指定插件的指定方法，并传入相应参数",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]
		methodName := args[1]
		methodArgs := args[2:]

		// 初始化彩色输出
		red := color.New(color.FgRed).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		blue := color.New(color.FgBlue).SprintFunc()

		// 加载插件
		pm := shared.NewPluginManager()
		defer pm.UnloadAll()

		client, ok := pm.Plugins[pluginName]
		if !ok {
			fmt.Println(red("错误: 插件未加载 - " + pluginName))
			return
		}

		// 获取插件实例
		rpcClient, err := client.Client()
		if err != nil {
			fmt.Println(red("RPC连接失败:"), err)
			return
		}

		raw, err := rpcClient.Dispense(pluginName)
		if err != nil {
			fmt.Println(red("获取插件实例失败:"), err)
			return
		}

		// 转换参数类型
		convertedArgs := make([]interface{}, len(methodArgs))
		for i, arg := range methodArgs {
			switch pluginName {
			case "calculator":
				if num, err := strconv.Atoi(arg); err == nil {
					convertedArgs[i] = num
				} else {
					convertedArgs[i] = arg
				}
			case "date_utils":
				if methodName == "AddDays" && i == 0 {
					if t, err := time.Parse("2006-01-02", arg); err == nil {
						convertedArgs[i] = t
					} else {
						convertedArgs[i] = arg
					}
				}
			default:
				convertedArgs[i] = arg
			}
		}

		// 调用方法
		var result interface{}
		switch p := raw.(type) {
		case shared.DynamicPlugin:
			result, err = p.Invoke(methodName, convertedArgs...)
		default:
			err = fmt.Errorf("不支持的插件类型")
		}

		if err != nil {
			fmt.Println(red("调用失败:"), err)
			return
		}

		fmt.Printf("%s %s.%s(%s)\n",
			blue("调用结果:"),
			green(pluginName),
			green(methodName),
			strings.Join(methodArgs, ", "))
		fmt.Printf("%s %v\n", blue("返回值:"), result)
	},
}

func init() {
	rootCmd.AddCommand(invokeCmd)
}
