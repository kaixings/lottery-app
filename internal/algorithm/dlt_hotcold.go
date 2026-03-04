package algorithm

import (
	"fmt"
	"sort"

	"lottery-app/internal/storage"
)

// DLTHotColdAlgo 冷热号算法：3个热号 + 2个冷号组成前区，后区取热号
type DLTHotColdAlgo struct {
	HotWindow int
}

func (h *DLTHotColdAlgo) Name() string { return "hot-cold" }

func (h *DLTHotColdAlgo) Predict(records []storage.DLTRecord) (*DLTResult, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("没有历史数据")
	}

	window := h.HotWindow
	if window <= 0 {
		window = 10
	}
	if window > len(records) {
		window = len(records)
	}

	// 前区热号
	hotFreq := make(map[int]int)
	for _, r := range records[:window] {
		for _, n := range r.Front {
			hotFreq[n]++
		}
	}

	// 前区冷号：统计距上次出现的间隔
	lastSeen := make(map[int]int)
	for i, r := range records {
		for _, n := range r.Front {
			if _, ok := lastSeen[n]; !ok {
				lastSeen[n] = i
			}
		}
	}

	type coldPair struct{ num, gap int }
	coldPairs := make([]coldPair, 0, 35)
	for i := 1; i <= 35; i++ {
		gap := len(records)
		if idx, ok := lastSeen[i]; ok {
			gap = idx
		}
		coldPairs = append(coldPairs, coldPair{i, gap})
	}
	sort.Slice(coldPairs, func(i, j int) bool {
		return coldPairs[i].gap > coldPairs[j].gap
	})

	// 取3个热号 + 2个冷号 = 5个前区
	selected := make(map[int]bool)
	front := make([]int, 0, 5)
	for _, n := range topN(hotFreq, 5, 1, 35) {
		if len(front) >= 3 {
			break
		}
		selected[n] = true
		front = append(front, n)
	}
	for _, cp := range coldPairs {
		if len(front) >= 5 {
			break
		}
		if !selected[cp.num] {
			front = append(front, cp.num)
		}
	}
	sort.Ints(front)

	// 后区取热号
	backFreq := make(map[int]int)
	for _, r := range records[:window] {
		for _, n := range r.Back {
			backFreq[n]++
		}
	}
	back := topN(backFreq, 2, 1, 12)

	return &DLTResult{
		Front: front,
		Back:  back,
		Desc:  fmt.Sprintf("热号(近%d期)3个 + 冷号2个", window),
	}, nil
}

func (h *DLTHotColdAlgo) PredictMultiple(records []storage.DLTRecord, count int) ([]*DLTResult, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("没有历史数据")
	}

	window := h.HotWindow
	if window <= 0 {
		window = 10
	}
	if window > len(records) {
		window = len(records)
	}

	hotFreq := make(map[int]int)
	for _, r := range records[:window] {
		for _, n := range r.Front {
			hotFreq[n]++
		}
	}

	lastSeen := make(map[int]int)
	for i, r := range records {
		for _, n := range r.Front {
			if _, ok := lastSeen[n]; !ok {
				lastSeen[n] = i
			}
		}
	}

	type coldPair struct{ num, gap int }
	coldPairs := make([]coldPair, 0, 35)
	for i := 1; i <= 35; i++ {
		gap := len(records)
		if idx, ok := lastSeen[i]; ok {
			gap = idx
		}
		coldPairs = append(coldPairs, coldPair{i, gap})
	}
	sort.Slice(coldPairs, func(i, j int) bool {
		return coldPairs[i].gap > coldPairs[j].gap
	})

	hotTop := topN(hotFreq, 35, 1, 35)

	backFreq := make(map[int]int)
	for _, r := range records[:window] {
		for _, n := range r.Back {
			backFreq[n]++
		}
	}
	backTop := topN(backFreq, 12, 1, 12)

	results := make([]*DLTResult, 0, count)
	for i := 0; i < count; i++ {
		selected := make(map[int]bool)
		front := make([]int, 0, 5)

		hotStart := (i * 2) % len(hotTop)
		coldStart := (i * 2) % len(coldPairs)

		for j := hotStart; len(front) < 3 && j < len(hotTop); j++ {
			n := hotTop[j]
			if !selected[n] {
				selected[n] = true
				front = append(front, n)
			}
		}
		for j := coldStart; len(front) < 5 && j < len(coldPairs); j++ {
			n := coldPairs[j].num
			if !selected[n] {
				front = append(front, n)
			}
		}
		sort.Ints(front)

		// 后区：2个不同后区号轮换
		b1 := (i * 2) % len(backTop)
		b2 := (i*2 + 1) % len(backTop)
		back := []int{backTop[b1], backTop[b2]}
		sort.Ints(back)

		results = append(results, &DLTResult{
			Front: front,
			Back:  back,
			Desc:  fmt.Sprintf("热号(近%d期)3个 + 冷号2个 (第%d组)", window, i+1),
		})
	}
	return results, nil
}
