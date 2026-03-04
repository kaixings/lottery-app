package algorithm

import "lottery-app/internal/storage"

// DLTResult 大乐透预测结果
type DLTResult struct {
	Front []int
	Back  []int
	Desc  string
}

// DLTAlgorithm 大乐透预测算法接口
type DLTAlgorithm interface {
	Name() string
	Predict(records []storage.DLTRecord) (*DLTResult, error)
	PredictMultiple(records []storage.DLTRecord, count int) ([]*DLTResult, error)
}

// DLTAll 返回所有可用的大乐透算法
func DLTAll() []DLTAlgorithm {
	return []DLTAlgorithm{
		&DLTFrequencyAlgo{},
		&DLTHotColdAlgo{},
		&DLTRandomAlgo{},
		&DLTWeightedAlgo{},
	}
}

// DLTGet 按名称获取大乐透算法
func DLTGet(name string) DLTAlgorithm {
	for _, a := range DLTAll() {
		if a.Name() == name {
			return a
		}
	}
	return nil
}
