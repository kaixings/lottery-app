package algorithm

import "lottery-app/internal/storage"

// Result 预测结果
type Result struct {
	Red  []int
	Blue int
	Desc string // 算法说明
}

// Algorithm 预测算法接口
type Algorithm interface {
	Name() string
	Predict(records []storage.LotteryRecord) (*Result, error)
	PredictMultiple(records []storage.LotteryRecord, count int) ([]*Result, error)
}

// All 返回所有可用算法
func All() []Algorithm {
	return []Algorithm{
		&FrequencyAlgo{},
		&HotColdAlgo{},
		&RandomAlgo{},
		&WeightedAlgo{},
	}
}

// Get 按名称获取算法
func Get(name string) Algorithm {
	for _, a := range All() {
		if a.Name() == name {
			return a
		}
	}
	return nil
}
