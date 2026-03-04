package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

// LotteryRecord 双色球开奖记录
type LotteryRecord struct {
	Issue string `json:"issue"` // 期号
	Date  string `json:"date"`  // 开奖日期
	Red   []int  `json:"red"`   // 红球 6个 1-33
	Blue  int    `json:"blue"`  // 蓝球 1-16
}

// Store 本地数据存储
type Store struct {
	filePath string
}

type storeData struct {
	Records []LotteryRecord `json:"records"`
}

func NewStore(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	return &Store{filePath: filepath.Join(dataDir, "lottery.json")}, nil
}

func (s *Store) Load() ([]LotteryRecord, error) {
	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return []LotteryRecord{}, nil
	}
	if err != nil {
		return nil, err
	}
	var sd storeData
	if err := json.Unmarshal(data, &sd); err != nil {
		return nil, err
	}
	return sd.Records, nil
}

func (s *Store) Save(records []LotteryRecord) error {
	// 按期号排序（降序）
	sort.Slice(records, func(i, j int) bool {
		return records[i].Issue > records[j].Issue
	})
	data, err := json.MarshalIndent(storeData{Records: records}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

// Merge 合并新记录（去重）
func (s *Store) Merge(existing, newRecords []LotteryRecord) []LotteryRecord {
	seen := make(map[string]bool)
	for _, r := range existing {
		seen[r.Issue] = true
	}
	result := append([]LotteryRecord{}, existing...)
	for _, r := range newRecords {
		if !seen[r.Issue] {
			result = append(result, r)
			seen[r.Issue] = true
		}
	}
	return result
}
