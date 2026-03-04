package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	label string
	view  viewState
}

type menuModel struct {
	lotteryName string
	items       []menuItem
	cursor      int
}

type switchViewMsg struct{ view viewState }

func newMenuModel(lotteryName string) menuModel {
	return menuModel{
		lotteryName: lotteryName,
		items: []menuItem{
			{"抓取数据", viewFetch},
			{"查看历史", viewHistory},
			{"预测号码", viewPredict},
			{"返回", viewLotterySelect},
			{"退出", -1},
		},
		cursor: 0,
	}
}

func (m menuModel) Init() tea.Cmd { return nil }

func (m menuModel) Update(msg tea.Msg) (menuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.items) - 1
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter":
			selected := m.items[m.cursor]
			if selected.view == -1 {
				return m, tea.Quit
			}
			return m, func() tea.Msg { return switchViewMsg{selected.view} }
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menuModel) View() string {
	s := titleStyle.Render(m.lotteryName+"预测工具") + "\n\n"
	for i, item := range m.items {
		if i == m.cursor {
			s += selectedStyle.Render("> "+item.label) + "\n"
		} else {
			s += normalStyle.Render(fmt.Sprintf("  %s", item.label)) + "\n"
		}
	}
	s += "\n" + dimStyle.Render("↑↓ 移动  Enter 确认  q 退出")
	return s
}
