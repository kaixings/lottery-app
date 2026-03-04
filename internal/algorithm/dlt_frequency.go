package algorithm

import (
	"fmt"
	"sort"

	"lottery-app/internal/storage"
)

// DLTFrequencyAlgo 频率算法：选取历史出现频率最高的号码
type DLTFrequencyAlgo struct{}

func (f *DLTFrequencyAlgo) Name() string { return "frequency" }

func (f *DLTFrequencyAlgo) Predict(records []storage.DLTRecord) (*DLTResult, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("没有历史数据")
	}

	frontFreq := make(map[int]int)
	backFreq := make(map[int]int)
	for _, r := range records {
		for _, n := range r.Front {
			frontFreq[n]++
		}
		for _, n := range r.Back {
			backFreq[n]++
		}
	}

	front := topN(frontFreq, 5, 1, 35)
	back := topN(backFreq, 2, 1, 12)

	return &DLTResult{
		Front: front,
		Back:  back,
		Desc:  fmt.Sprintf("基于 %d 期历史数据的出现频率", len(records)),
	}, nil
}

func (f *DLTFrequencyAlgo) PredictMultiple(records []storage.DLTRecord, count int) ([]*DLTResult, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("没有历史数据")
	}

	frontFreq := make(map[int]int)
	backFreq := make(map[int]int)
	for _, r := range records {
		for _, n := range r.Front {
			frontFreq[n]++
		}
		for _, n := range r.Back {
			backFreq[n]++
		}
	}

	frontSorted := sortedByFreq(frontFreq, 1, 35) // 35个，索引 0-34
	backSorted := sortedByFreq(backFreq, 1, 12)   // 12个，索引 0-11

	results := make([]*DLTResult, 0, count)
	for i := 0; i < count; i++ {
		start := i * 5
		if start+5 > len(frontSorted) {
			start = start % (len(frontSorted) - 4)
		}
		front := make([]int, 5)
		copy(front, frontSorted[start:start+5])
		sort.Ints(front)

		// 后区：取2个不同的后区号，轮换偏移
		b1 := i % (len(backSorted) - 1)
		back := []int{backSorted[b1], backSorted[b1+1]}
		sort.Ints(back)

		results = append(results, &DLTResult{
			Front: front,
			Back:  back,
			Desc:  fmt.Sprintf("基于 %d 期历史数据的出现频率 (第%d组)", len(records), i+1),
		})
	}
	return results, nil
}
