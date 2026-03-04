package algorithm

import (
	"fmt"
	"sort"

	"lottery-app/internal/storage"
)

// HotColdAlgo 冷热号算法：结合近期热号和长期未出现的冷号
type HotColdAlgo struct {
	HotWindow int // 热号统计窗口（最近N期）
}

func (h *HotColdAlgo) Name() string { return "hot-cold" }

func (h *HotColdAlgo) Predict(records []storage.LotteryRecord) (*Result, error) {
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

	// 热号：最近window期出现的红球
	hotFreq := make(map[int]int)
	for _, r := range records[:window] {
		for _, n := range r.Red {
			hotFreq[n]++
		}
	}

	// 冷号：统计每个号码距上次出现的间隔
	lastSeen := make(map[int]int) // 号码 -> 最近出现的期数索引
	for i, r := range records {
		for _, n := range r.Red {
			if _, ok := lastSeen[n]; !ok {
				lastSeen[n] = i
			}
		}
	}

	// 冷号 = 间隔最大的号码
	type coldPair struct {
		num, gap int
	}
	coldPairs := make([]coldPair, 0, 33)
	for i := 1; i <= 33; i++ {
		gap := len(records) // 从未出现
		if idx, ok := lastSeen[i]; ok {
			gap = idx
		}
		coldPairs = append(coldPairs, coldPair{i, gap})
	}
	sort.Slice(coldPairs, func(i, j int) bool {
		return coldPairs[i].gap > coldPairs[j].gap
	})

	// 取3个热号 + 3个冷号
	selected := make(map[int]bool)
	red := make([]int, 0, 6)

	hotTop := topN(hotFreq, 6, 1, 33)
	for _, n := range hotTop {
		if len(red) >= 3 {
			break
		}
		selected[n] = true
		red = append(red, n)
	}
	for _, cp := range coldPairs {
		if len(red) >= 6 {
			break
		}
		if !selected[cp.num] {
			red = append(red, cp.num)
		}
	}
	sort.Ints(red)

	// 蓝球取热号
	blueFreq := make(map[int]int)
	for _, r := range records[:window] {
		blueFreq[r.Blue]++
	}
	blue := topN(blueFreq, 1, 1, 16)[0]

	return &Result{
		Red:  red,
		Blue: blue,
		Desc: fmt.Sprintf("热号(近%d期)3个 + 冷号3个", window),
	}, nil
}

func (h *HotColdAlgo) PredictMultiple(records []storage.LotteryRecord, count int) ([]*Result, error) {
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

	// 热号：最近window期出现的红球
	hotFreq := make(map[int]int)
	for _, r := range records[:window] {
		for _, n := range r.Red {
			hotFreq[n]++
		}
	}

	// 冷号：统计每个号码距上次出现的间隔
	lastSeen := make(map[int]int)
	for i, r := range records {
		for _, n := range r.Red {
			if _, ok := lastSeen[n]; !ok {
				lastSeen[n] = i
			}
		}
	}

	// 冷号 = 间隔最大的号码
	type coldPair struct {
		num, gap int
	}
	coldPairs := make([]coldPair, 0, 33)
	for i := 1; i <= 33; i++ {
		gap := len(records)
		if idx, ok := lastSeen[i]; ok {
			gap = idx
		}
		coldPairs = append(coldPairs, coldPair{i, gap})
	}
	sort.Slice(coldPairs, func(i, j int) bool {
		return coldPairs[i].gap > coldPairs[j].gap
	})

	hotTop := topN(hotFreq, 33, 1, 33) // 获取所有热号排序

	// 蓝球热号列表
	blueFreq := make(map[int]int)
	for _, r := range records[:window] {
		blueFreq[r.Blue]++
	}
	blueTop := topN(blueFreq, 16, 1, 16)

	results := make([]*Result, 0, count)
	for i := 0; i < count; i++ {
		selected := make(map[int]bool)
		red := make([]int, 0, 6)

		// 每组使用不同的热号和冷号组合
		hotStart := (i * 2) % len(hotTop)
		coldStart := (i * 2) % len(coldPairs)

		// 取3个热号
		for j := hotStart; len(red) < 3 && j < len(hotTop); j++ {
			n := hotTop[j]
			if !selected[n] {
				selected[n] = true
				red = append(red, n)
			}
		}

		// 取3个冷号
		for j := coldStart; len(red) < 6 && j < len(coldPairs); j++ {
			n := coldPairs[j].num
			if !selected[n] {
				red = append(red, n)
			}
		}
		sort.Ints(red)

		// 蓝球轮换
		blueIdx := i % len(blueTop)
		blue := blueTop[blueIdx]

		results = append(results, &Result{
			Red:  red,
			Blue: blue,
			Desc: fmt.Sprintf("热号(近%d期)3个 + 冷号3个 (第%d组)", window, i+1),
		})
	}

	return results, nil
}
