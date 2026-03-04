package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"lottery-app/internal/algorithm"
	"lottery-app/internal/storage"
)

type predictModel struct {
	algos     []algorithm.Algorithm
	algoIdx   int
	count     int
	results   []algorithm.Result
	dataDir   string
	err       error
}

type predictResultMsg struct {
	results []algorithm.Result
	err     error
}

func newPredictModel(dataDir string) predictModel {
	return predictModel{
		algos:   algorithm.All(),
		algoIdx: 0,
		count:   1,
		dataDir: dataDir,
	}
}

func (m predictModel) generateCmd() tea.Cmd {
	return func() tea.Msg {
		store, err := storage.NewStore(m.dataDir)
		if err != nil {
			return predictResultMsg{err: err}
		}
		records, err := store.Load()
		if err != nil {
			return predictResultMsg{err: err}
		}
		algo := m.algos[m.algoIdx]

		// 使用 PredictMultiple 方法生成多组预测
		resultPtrs, err := algo.PredictMultiple(records, m.count)
		if err != nil {
			return predictResultMsg{err: err}
		}

		// 转换为非指针类型
		results := make([]algorithm.Result, len(resultPtrs))
		for i, r := range resultPtrs {
			results[i] = *r
		}

		return predictResultMsg{results: results}
	}
}

func (m predictModel) Init() tea.Cmd { return nil }

func (m predictModel) Update(msg tea.Msg) (predictModel, tea.Cmd) {
	switch msg := msg.(type) {
	case predictResultMsg:
		m.results = msg.results
		m.err = msg.err
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.algoIdx > 0 {
				m.algoIdx--
			} else {
				m.algoIdx = len(m.algos) - 1
			}
			m.results = nil
		case "down", "j":
			if m.algoIdx < len(m.algos)-1 {
				m.algoIdx++
			} else {
				m.algoIdx = 0
			}
			m.results = nil
		case "+", "right", "l":
			if m.count < 10 {
				m.count++
			}
			m.results = nil
		case "-", "left", "h":
			if m.count > 1 {
				m.count--
			}
			m.results = nil
		case "enter":
			return m, m.generateCmd()
		case "esc", "q":
			return m, func() tea.Msg { return switchViewMsg{viewMenu} }
		}
	}
	return m, nil
}

func (m predictModel) View() string {
	s := titleStyle.Render("预测号码") + "\n\n"

	s += tableHeaderStyle.Render("选择算法：") + "\n"
	for i, a := range m.algos {
		if i == m.algoIdx {
			s += selectedStyle.Render("> "+a.Name()) + "\n"
		} else {
			s += normalStyle.Render("  "+a.Name()) + "\n"
		}
	}

	s += "\n" + normalStyle.Render(fmt.Sprintf("注数：%d", m.count)) + "\n"
	s += dimStyle.Render("+/→ 增加  -/← 减少") + "\n\n"

	if m.err != nil {
		s += errorStyle.Render("预测失败: "+m.err.Error()) + "\n"
	} else if len(m.results) == 0 {
		s += dimStyle.Render("按 Enter 生成预测") + "\n"
	} else {
		s += tableHeaderStyle.Render("预测结果：") + "\n"
		for i, r := range m.results {
			reds := make([]string, len(r.Red))
			for j, v := range r.Red {
				reds[j] = fmt.Sprintf("%02d", v)
			}
			line := fmt.Sprintf("  第%d注: 红球 %s  蓝球 %02d", i+1, strings.Join(reds, " "), r.Blue)
			s += tableRowStyle.Render(line) + "\n"
		}
	}

	s += "\n" + dimStyle.Render("↑↓ 切换算法  Enter 生成  Esc/q 返回")
	return s
}
