package algorithm

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"lottery-app/internal/storage"
)

// DLTRandomAlgo 纯随机算法
type DLTRandomAlgo struct{}

func (r *DLTRandomAlgo) Name() string { return "random" }

func (r *DLTRandomAlgo) Predict(_ []storage.DLTRecord) (*DLTResult, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	pool := make([]int, 35)
	for i := range pool {
		pool[i] = i + 1
	}
	rng.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	front := append([]int{}, pool[:5]...)
	sort.Ints(front)

	backPool := make([]int, 12)
	for i := range backPool {
		backPool[i] = i + 1
	}
	rng.Shuffle(len(backPool), func(i, j int) { backPool[i], backPool[j] = backPool[j], backPool[i] })
	back := append([]int{}, backPool[:2]...)
	sort.Ints(back)

	return &DLTResult{Front: front, Back: back, Desc: "纯随机选号"}, nil
}

func (r *DLTRandomAlgo) PredictMultiple(_ []storage.DLTRecord, count int) ([]*DLTResult, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	results := make([]*DLTResult, 0, count)

	for i := 0; i < count; i++ {
		pool := make([]int, 35)
		for j := range pool {
			pool[j] = j + 1
		}
		rng.Shuffle(len(pool), func(j, k int) { pool[j], pool[k] = pool[k], pool[j] })
		front := append([]int{}, pool[:5]...)
		sort.Ints(front)

		backPool := make([]int, 12)
		for j := range backPool {
			backPool[j] = j + 1
		}
		rng.Shuffle(len(backPool), func(j, k int) { backPool[j], backPool[k] = backPool[k], backPool[j] })
		back := append([]int{}, backPool[:2]...)
		sort.Ints(back)

		results = append(results, &DLTResult{
			Front: front,
			Back:  back,
			Desc:  fmt.Sprintf("纯随机选号 (第%d组)", i+1),
		})
	}
	return results, nil
}

// DLTWeightedAlgo 加权随机算法：按历史频率加权随机
type DLTWeightedAlgo struct{}

func (w *DLTWeightedAlgo) Name() string { return "weighted" }

func (w *DLTWeightedAlgo) Predict(records []storage.DLTRecord) (*DLTResult, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("没有历史数据")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	frontFreq := make([]int, 36) // index 1-35
	backFreq := make([]int, 13)  // index 1-12
	for _, r := range records {
		for _, n := range r.Front {
			frontFreq[n]++
		}
		for _, n := range r.Back {
			backFreq[n]++
		}
	}

	front := weightedSample(rng, frontFreq, 1, 35, 5)
	back := weightedSample(rng, backFreq, 1, 12, 2)

	return &DLTResult{
		Front: front,
		Back:  back,
		Desc:  fmt.Sprintf("基于 %d 期历史数据的加权随机", len(records)),
	}, nil
}

func (w *DLTWeightedAlgo) PredictMultiple(records []storage.DLTRecord, count int) ([]*DLTResult, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("没有历史数据")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	frontFreq := make([]int, 36)
	backFreq := make([]int, 13)
	for _, r := range records {
		for _, n := range r.Front {
			frontFreq[n]++
		}
		for _, n := range r.Back {
			backFreq[n]++
		}
	}

	results := make([]*DLTResult, 0, count)
	for i := 0; i < count; i++ {
		front := weightedSample(rng, frontFreq, 1, 35, 5)
		back := weightedSample(rng, backFreq, 1, 12, 2)
		results = append(results, &DLTResult{
			Front: front,
			Back:  back,
			Desc:  fmt.Sprintf("基于 %d 期历史数据的加权随机 (第%d组)", len(records), i+1),
		})
	}
	return results, nil
}
