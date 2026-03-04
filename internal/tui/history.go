package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"lottery-app/internal/storage"
)

const historyVisibleRows = 15

type historyModel struct {
	records []storage.LotteryRecord
	offset  int
	dataDir string
	err     error
}

type historyLoadedMsg struct {
	records []storage.LotteryRecord
	err     error
}

func newHistoryModel(dataDir string) historyModel {
	return historyModel{dataDir: dataDir}
}

func (m historyModel) loadCmd() tea.Cmd {
	return func() tea.Msg {
		store, err := storage.NewStore(m.dataDir)
		if err != nil {
			return historyLoadedMsg{err: err}
		}
		records, err := store.Load()
		return historyLoadedMsg{records: records, err: err}
	}
}

func (m historyModel) Init() tea.Cmd { return m.loadCmd() }

func (m historyModel) Update(msg tea.Msg) (historyModel, tea.Cmd) {
	switch msg := msg.(type) {
	case historyLoadedMsg:
		m.records = msg.records
		m.err = msg.err
		m.offset = 0
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.offset > 0 {
				m.offset--
			}
		case "down", "j":
			if m.offset < len(m.records)-historyVisibleRows {
				m.offset++
			}
		case "pgup", "u":
			m.offset -= historyVisibleRows
			if m.offset < 0 {
				m.offset = 0
			}
		case "pgdown", "d":
			m.offset += historyVisibleRows
			max := len(m.records) - historyVisibleRows
			if max < 0 {
				max = 0
			}
			if m.offset > max {
				m.offset = max
			}
		case "esc", "q":
			return m, func() tea.Msg { return switchViewMsg{viewMenu} }
		}
	}
	return m, nil
}

func (m historyModel) View() string {
	s := titleStyle.Render("历史开奖记录") + "\n\n"
	if m.err != nil {
		return s + errorStyle.Render("加载失败: "+m.err.Error())
	}
	if len(m.records) == 0 {
		return s + dimStyle.Render("暂无数据，请先抓取数据") + "\n\n" + dimStyle.Render("Esc/q 返回")
	}

	header := fmt.Sprintf("%-12s  %-12s  %-20s  %s", "期号", "日期", "红球", "蓝球")
	s += tableHeaderStyle.Render(header) + "\n"
	s += tableHeaderStyle.Render(strings.Repeat("─", 60)) + "\n"

	end := m.offset + historyVisibleRows
	if end > len(m.records) {
		end = len(m.records)
	}
	for _, r := range m.records[m.offset:end] {
		reds := make([]string, len(r.Red))
		for i, v := range r.Red {
			reds[i] = fmt.Sprintf("%02d", v)
		}
		line := fmt.Sprintf("%-12s  %-12s  %-20s  %02d",
			r.Issue, r.Date, strings.Join(reds, " "), r.Blue)
		s += tableRowStyle.Render(line) + "\n"
	}

	s += "\n" + dimStyle.Render(fmt.Sprintf("%d/%d  ↑↓ 滚动  PgUp/u PgDn/d 翻页  Esc/q 返回", m.offset+1, len(m.records)))
	return s
}
