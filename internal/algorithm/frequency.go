package algorithm

import (
	"fmt"
	"sort"

	"lottery-app/internal/storage"
)

// FrequencyAlgo 频率算法：选取历史出现频率最高的号码
type FrequencyAlgo struct{}

func (f *FrequencyAlgo) Name() string { return "frequency" }

func (f *FrequencyAlgo) Predict(records []storage.LotteryRecord) (*Result, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("没有历史数据")
	}

	redFreq := make(map[int]int)
	blueFreq := make(map[int]int)

	for _, r := range records {
		for _, n := range r.Red {
			redFreq[n]++
		}
		blueFreq[r.Blue]++
	}

	red := topN(redFreq, 6, 1, 33)
	blue := topN(blueFreq, 1, 1, 16)[0]

	return &Result{
		Red:  red,
		Blue: blue,
		Desc: fmt.Sprintf("基于 %d 期历史数据的出现频率", len(records)),
	}, nil
}

func (f *FrequencyAlgo) PredictMultiple(records []storage.LotteryRecord, count int) ([]*Result, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("没有历史数据")
	}

	redFreq := make(map[int]int)
	blueFreq := make(map[int]int)

	for _, r := range records {
		for _, n := range r.Red {
			redFreq[n]++
		}
		blueFreq[r.Blue]++
	}

	// 获取排序后的红球和蓝球列表
	redSorted := sortedByFreq(redFreq, 1, 33)
	blueSorted := sortedByFreq(blueFreq, 1, 16)

	results := make([]*Result, 0, count)
	for i := 0; i < count; i++ {
		// 每组取6个红球，从不同的起始位置
		start := i * 6
		if start+6 > len(redSorted) {
			// 如果号码不够，循环使用
			start = start % (len(redSorted) - 5)
		}

		red := make([]int, 6)
		copy(red, redSorted[start:start+6])
		sort.Ints(red)

		// 蓝球也轮换
		blueIdx := i % len(blueSorted)
		blue := blueSorted[blueIdx]

		results = append(results, &Result{
			Red:  red,
			Blue: blue,
			Desc: fmt.Sprintf("基于 %d 期历史数据的出现频率 (第%d组)", len(records), i+1),
		})
	}

	return results, nil
}

// sortedByFreq 返回按频率排序的号码列表
func sortedByFreq(freq map[int]int, min, max int) []int {
	type pair struct {
		num, cnt int
	}
	pairs := make([]pair, 0, max-min+1)
	for i := min; i <= max; i++ {
		pairs = append(pairs, pair{i, freq[i]})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].cnt != pairs[j].cnt {
			return pairs[i].cnt > pairs[j].cnt
		}
		return pairs[i].num < pairs[j].num
	})
	result := make([]int, len(pairs))
	for i := range pairs {
		result[i] = pairs[i].num
	}
	return result
}

// topN 从频率map中取出频率最高的N个号码（范围 min-max）
func topN(freq map[int]int, n, min, max int) []int {
	type pair struct {
		num, cnt int
	}
	pairs := make([]pair, 0, max-min+1)
	for i := min; i <= max; i++ {
		pairs = append(pairs, pair{i, freq[i]})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].cnt != pairs[j].cnt {
			return pairs[i].cnt > pairs[j].cnt
		}
		return pairs[i].num < pairs[j].num
	})
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = pairs[i].num
	}
	sort.Ints(result)
	return result
}
