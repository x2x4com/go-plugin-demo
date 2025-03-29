package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "plugin-cli",
	Short: "Go插件系统命令行工具",
	Long: `Go插件系统命令行工具，支持插件管理和方法调用。

支持以下子命令:
  list    - 列出所有插件及方法
  invoke  - 调用插件方法
  help    - 显示详细帮助信息`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 初始化彩色输出
	color.NoColor = false
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		// 自定义帮助输出颜色
		blue := color.New(color.FgBlue).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Printf("%s\n\n", blue(cmd.Long))
		fmt.Printf("%s\n", green("可用命令:"))
		for _, cmd := range cmd.Commands() {
			fmt.Printf("  %-15s %s\n", cmd.Name(), cmd.Short)
		}
		fmt.Printf("\n使用 \"%s [command] --help\" 查看命令详情\n", cmd.CommandPath())
	})
}
