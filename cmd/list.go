package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"lottery-app/internal/storage"
)

var listCount int

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "查看历史开奖数据",
	Example: `  lottery list           # 显示最近10期
  lottery list -n 20     # 显示最近20期`,
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
			fmt.Println("本地没有历史数据，请先运行: lottery fetch")
			return nil
		}

		n := listCount
		if n > len(records) {
			n = len(records)
		}

		fmt.Printf("%-12s  %-12s  %-25s  %s\n", "期号", "日期", "红球", "蓝球")
		fmt.Println(strings.Repeat("-", 65))
		for _, r := range records[:n] {
			redStrs := make([]string, len(r.Red))
			for i, n := range r.Red {
				redStrs[i] = fmt.Sprintf("%02d", n)
			}
			fmt.Printf("%-12s  %-12s  %-25s  %02d\n",
				r.Issue, r.Date, strings.Join(redStrs, " "), r.Blue)
		}
		fmt.Printf("\n共 %d 期数据\n", len(records))
		return nil
	},
}

func init() {
	listCmd.Flags().IntVarP(&listCount, "count", "n", 10, "显示期数")
	rootCmd.AddCommand(listCmd)
}
