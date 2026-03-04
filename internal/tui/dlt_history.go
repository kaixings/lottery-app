package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"lottery-app/internal/storage"
)

type dltHistoryModel struct {
	records []storage.DLTRecord
	offset  int
	dataDir string
	err     error
}

type dltHistoryLoadedMsg struct {
	records []storage.DLTRecord
	err     error
}

func newDLTHistoryModel(dataDir string) dltHistoryModel {
	return dltHistoryModel{dataDir: dataDir}
}

func (m dltHistoryModel) loadCmd() tea.Cmd {
	return func() tea.Msg {
		store, err := storage.NewDLTStore(m.dataDir)
		if err != nil {
			return dltHistoryLoadedMsg{err: err}
		}
		records, err := store.Load()
		return dltHistoryLoadedMsg{records: records, err: err}
	}
}

func (m dltHistoryModel) Init() tea.Cmd { return m.loadCmd() }

func (m dltHistoryModel) Update(msg tea.Msg) (dltHistoryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case dltHistoryLoadedMsg:
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

func (m dltHistoryModel) View() string {
	s := titleStyle.Render("历史开奖记录（大乐透）") + "\n\n"
	if m.err != nil {
		return s + errorStyle.Render("加载失败: "+m.err.Error())
	}
	if len(m.records) == 0 {
		return s + dimStyle.Render("暂无数据，请先抓取数据") + "\n\n" + dimStyle.Render("Esc/q 返回")
	}

	header := fmt.Sprintf("%-12s  %-12s  %-18s  %s", "期号", "日期", "前区", "后区")
	s += tableHeaderStyle.Render(header) + "\n"
	s += tableHeaderStyle.Render(strings.Repeat("─", 60)) + "\n"

	end := m.offset + historyVisibleRows
	if end > len(m.records) {
		end = len(m.records)
	}
	for _, r := range m.records[m.offset:end] {
		fronts := make([]string, len(r.Front))
		for i, v := range r.Front {
			fronts[i] = fmt.Sprintf("%02d", v)
		}
		backs := make([]string, len(r.Back))
		for i, v := range r.Back {
			backs[i] = fmt.Sprintf("%02d", v)
		}
		line := fmt.Sprintf("%-12s  %-12s  %-18s  %s",
			r.Issue, r.Date, strings.Join(fronts, " "), strings.Join(backs, " "))
		s += tableRowStyle.Render(line) + "\n"
	}

	s += "\n" + dimStyle.Render(fmt.Sprintf("%d/%d  ↑↓ 滚动  PgUp/u PgDn/d 翻页  Esc/q 返回",
		m.offset+1, len(m.records)))
	return s
}
