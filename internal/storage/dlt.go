package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

// DLTRecord 大乐透开奖记录
type DLTRecord struct {
	Issue string `json:"issue"` // 期号
	Date  string `json:"date"`  // 开奖日期
	Front []int  `json:"front"` // 前区 5个 1-35
	Back  []int  `json:"back"`  // 后区 2个 1-12
}

// DLTStore 大乐透本地数据存储
type DLTStore struct {
	filePath string
}

type dltStoreData struct {
	Records []DLTRecord `json:"records"`
}

func NewDLTStore(dataDir string) (*DLTStore, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	return &DLTStore{filePath: filepath.Join(dataDir, "dlt.json")}, nil
}

func (s *DLTStore) Load() ([]DLTRecord, error) {
	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return []DLTRecord{}, nil
	}
	if err != nil {
		return nil, err
	}
	var sd dltStoreData
	if err := json.Unmarshal(data, &sd); err != nil {
		return nil, err
	}
	return sd.Records, nil
}

func (s *DLTStore) Save(records []DLTRecord) error {
	sort.Slice(records, func(i, j int) bool {
		return records[i].Issue > records[j].Issue
	})
	data, err := json.MarshalIndent(dltStoreData{Records: records}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

func (s *DLTStore) Merge(existing, newRecords []DLTRecord) []DLTRecord {
	seen := make(map[string]bool)
	for _, r := range existing {
		seen[r.Issue] = true
	}
	result := append([]DLTRecord{}, existing...)
	for _, r := range newRecords {
		if !seen[r.Issue] {
			result = append(result, r)
			seen[r.Issue] = true
		}
	}
	return result
}
