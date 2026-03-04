package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type lotteryType int

const (
	lotterySSQ lotteryType = iota
	lotteryDLT
)

type lotterySelectModel struct {
	cursor int
}

type lotterySelectedMsg struct{ t lotteryType }

func newLotterySelectModel() lotterySelectModel {
	return lotterySelectModel{}
}

func (m lotterySelectModel) Init() tea.Cmd { return nil }

func (m lotterySelectModel) Update(msg tea.Msg) (lotterySelectModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = 1
			}
		case "down", "j":
			if m.cursor < 1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter":
			lt := lotteryType(m.cursor)
			return m, func() tea.Msg { return lotterySelectedMsg{lt} }
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m lotterySelectModel) View() string {
	s := titleStyle.Render("彩票预测工具") + "\n\n"
	s += normalStyle.Render("请选择彩种：") + "\n\n"
	options := []string{"双色球", "大乐透"}
	for i, opt := range options {
		if i == m.cursor {
			s += selectedStyle.Render("> "+opt) + "\n"
		} else {
			s += normalStyle.Render("  "+opt) + "\n"
		}
	}
	s += "\n" + dimStyle.Render("↑↓ 移动  Enter 确认  q 退出")
	return s
}
