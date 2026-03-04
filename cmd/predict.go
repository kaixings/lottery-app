package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"lottery-app/internal/algorithm"
	"lottery-app/internal/storage"
)

var (
	algoName string
	times    int
)

var predictCmd = &cobra.Command{
	Use:   "predict",
	Short: "预测下一期号码",
	Example: `  lottery predict                      # 使用所有算法预测
  lottery predict -a frequency         # 使用频率算法
  lottery predict -a hot-cold          # 使用冷热号算法
  lottery predict -a weighted -t 5     # 加权随机生成5注`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewStore(dataDir)
		if err != nil {
			return err
		}
		records, err := store.Load()
		if err != nil {
			return err
		}
		if len(records) == 0 {
			return fmt.Errorf("本地没有历史数据，请先运行: lottery fetch")
		}

		var algos []algorithm.Algorithm
		if algoName == "" {
			algos = algorithm.All()
		} else {
			a := algorithm.Get(algoName)
			if a == nil {
				return fmt.Errorf("未知算法: %s\n可用算法: frequency, hot-cold, random, weighted", algoName)
			}
			algos = []algorithm.Algorithm{a}
		}

		fmt.Printf("基于 %d 期历史数据预测\n\n", len(records))

		for _, a := range algos {
			fmt.Printf("【%s】%s\n", a.Name(), getAlgoDesc(a.Name()))
			for i := 0; i < times; i++ {
				result, err := a.Predict(records)
				if err != nil {
					fmt.Printf("  错误: %v\n", err)
					continue
				}
				redStrs := make([]string, len(result.Red))
				for j, n := range result.Red {
					redStrs[j] = fmt.Sprintf("%02d", n)
				}
				fmt.Printf("  红球: %s  蓝球: %02d\n",
					strings.Join(redStrs, " "), result.Blue)
			}
			fmt.Println()
		}
		return nil
	},
}

func getAlgoDesc(name string) string {
	switch name {
	case "frequency":
		return "频率最高的号码"
	case "hot-cold":
		return "冷热号组合"
	case "random":
		return "纯随机"
	case "weighted":
		return "加权随机"
	}
	return ""
}

func init() {
	predictCmd.Flags().StringVarP(&algoName, "algo", "a", "", "算法名称 (frequency/hot-cold/random/weighted)")
	predictCmd.Flags().IntVarP(&times, "times", "t", 1, "生成注数")
	rootCmd.AddCommand(predictCmd)
}
