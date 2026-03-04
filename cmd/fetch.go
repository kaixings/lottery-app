package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"lottery-app/internal/fetcher"
	"lottery-app/internal/storage"
)

var fetchCount int

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "抓取历史开奖数据",
	Example: `  lottery fetch           # 抓取最近30期
  lottery fetch -n 100    # 抓取最近100期`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewStore(dataDir)
		if err != nil {
			return fmt.Errorf("初始化存储失败: %w", err)
		}

		fmt.Printf("正在抓取最近 %d 期数据...\n", fetchCount)
		newRecords, err := fetcher.Fetch(fetchCount)
		if err != nil {
			return fmt.Errorf("抓取失败: %w", err)
		}

		existing, err := store.Load()
		if err != nil {
			return fmt.Errorf("加载本地数据失败: %w", err)
		}

		merged := store.Merge(existing, newRecords)
		if err := store.Save(merged); err != nil {
			return fmt.Errorf("保存失败: %w", err)
		}

		added := len(merged) - len(existing)
		fmt.Printf("抓取完成: 新增 %d 期，本地共 %d 期数据\n", added, len(merged))
		return nil
	},
}

func init() {
	fetchCmd.Flags().IntVarP(&fetchCount, "count", "n", 30, "抓取期数")
	rootCmd.AddCommand(fetchCmd)
}
