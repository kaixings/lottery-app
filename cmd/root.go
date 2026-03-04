package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"lottery-app/internal/tui"
)

var dataDir string

var rootCmd = &cobra.Command{
	Use:   "lottery",
	Short: "双色球号码预测工具",
	Long: `双色球号码预测工具

支持从官方网站抓取历史数据，并使用多种算法预测下一期号码。

算法列表:
  frequency  频率算法 - 选取历史出现频率最高的号码
  hot-cold   冷热号算法 - 结合近期热号和长期未出现的冷号
  random     随机算法 - 纯随机选号
  weighted   加权随机 - 按历史频率加权随机选号`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.Run(dataDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dataDir, "data", "./data", "数据存储目录")
}
