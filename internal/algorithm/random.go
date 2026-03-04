package algorithm

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"lottery-app/internal/storage"
)

// RandomAlgo 纯随机算法
type RandomAlgo struct{}

func (r *RandomAlgo) Name() string { return "random" }

func (r *RandomAlgo) Predict(_ []storage.LotteryRecord) (*Result, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	pool := make([]int, 33)
	for i := range pool {
		pool[i] = i + 1
	}
	rng.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })

	red := pool[:6]
	sort.Ints(red)
	blue := rng.Intn(16) + 1

	return &Result{
		Red:  red,
		Blue: blue,
		Desc: "纯随机选号",
	}, nil
}

func (r *RandomAlgo) PredictMultiple(_ []storage.LotteryRecord, count int) ([]*Result, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	results := make([]*Result, 0, count)

	for i := 0; i < count; i++ {
		pool := make([]int, 33)
		for j := range pool {
			pool[j] = j + 1
		}
		rng.Shuffle(len(pool), func(j, k int) { pool[j], pool[k] = pool[k], pool[j] })

		red := pool[:6]
		sort.Ints(red)
		blue := rng.Intn(16) + 1

		results = append(results, &Result{
			Red:  red,
			Blue: blue,
			Desc: fmt.Sprintf("纯随机选号 (第%d组)", i+1),
		})
	}

	return results, nil
}

// WeightedAlgo 加权随机算法：按历史频率加权随机
type WeightedAlgo struct{}

func (w *WeightedAlgo) Name() string { return "weighted" }

func (w *WeightedAlgo) Predict(records []storage.LotteryRecord) (*Result, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("没有历史数据")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	redFreq := make([]int, 34) // index 1-33
	blueFreq := make([]int, 17)
	for _, r := range records {
		for _, n := range r.Red {
			redFreq[n]++
		}
		blueFreq[r.Blue]++
	}

	red := weightedSample(rng, redFreq, 1, 33, 6)
	blue := weightedSample(rng, blueFreq, 1, 16, 1)[0]

	return &Result{
		Red:  red,
		Blue: blue,
		Desc: fmt.Sprintf("基于 %d 期历史数据的加权随机", len(records)),
	}, nil
}

func (w *WeightedAlgo) PredictMultiple(records []storage.LotteryRecord, count int) ([]*Result, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("没有历史数据")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	redFreq := make([]int, 34)
	blueFreq := make([]int, 17)
	for _, r := range records {
		for _, n := range r.Red {
			redFreq[n]++
		}
		blueFreq[r.Blue]++
	}

	results := make([]*Result, 0, count)
	for i := 0; i < count; i++ {
		red := weightedSample(rng, redFreq, 1, 33, 6)
		blue := weightedSample(rng, blueFreq, 1, 16, 1)[0]

		results = append(results, &Result{
			Red:  red,
			Blue: blue,
			Desc: fmt.Sprintf("基于 %d 期历史数据的加权随机 (第%d组)", len(records), i+1),
		})
	}

	return results, nil
}

// weightedSample 按权重无放回抽样
func weightedSample(rng *rand.Rand, freq []int, min, max, n int) []int {
	weights := make([]int, max-min+1)
	nums := make([]int, max-min+1)
	for i := range weights {
		nums[i] = min + i
		w := freq[min+i]
		if w == 0 {
			w = 1 // 保证每个号码都有机会
		}
		weights[i] = w
	}

	result := make([]int, 0, n)
	used := make([]bool, len(weights))

	for len(result) < n {
		total := 0
		for i, w := range weights {
			if !used[i] {
				total += w
			}
		}
		pick := rng.Intn(total)
		cum := 0
		for i, w := range weights {
			if used[i] {
				continue
			}
			cum += w
			if pick < cum {
				result = append(result, nums[i])
				used[i] = true
				break
			}
		}
	}
	sort.Ints(result)
	return result
}
