package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type viewState int

const (
	viewLotterySelect viewState = iota
	viewMenu
	viewHistory
	viewPredict
	viewFetch
	viewDLTHistory
	viewDLTPredict
	viewDLTFetch
)

// Model 顶层模型，管理视图切换
type Model struct {
	currentView   viewState
	lottery       lotteryType
	dataDir       string
	lotterySelect lotterySelectModel
	menu          menuModel
	history       historyModel
	predict       predictModel
	fetch         fetchModel
	dltHistory    dltHistoryModel
	dltPredict    dltPredictModel
	dltFetch      dltFetchModel
}

func newModel(dataDir string) Model {
	return Model{
		currentView:   viewLotterySelect,
		dataDir:       dataDir,
		lotterySelect: newLotterySelectModel(),
		menu:          newMenuModel(""),
		history:       newHistoryModel(dataDir),
		predict:       newPredictModel(dataDir),
		fetch:         newFetchModel(dataDir),
		dltHistory:    newDLTHistoryModel(dataDir),
		dltPredict:    newDLTPredictModel(dataDir),
		dltFetch:      newDLTFetchModel(dataDir),
	}
}

func (m Model) Init() tea.Cmd {
	return m.lotterySelect.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case lotterySelectedMsg:
		m.lottery = msg.t
		name := "双色球"
		if m.lottery == lotteryDLT {
			name = "大乐透"
		}
		m.menu = newMenuModel(name)
		m.currentView = viewMenu
		return m, nil

	case switchViewMsg:
		target := msg.view
		if target == viewLotterySelect {
			m.currentView = viewLotterySelect
			return m, nil
		}
		if target == viewMenu {
			m.currentView = viewMenu
			return m, nil
		}
		// 根据当前彩种路由到对应视图
		if m.lottery == lotteryDLT {
			switch target {
			case viewFetch:
				m.dltFetch = newDLTFetchModel(m.dataDir)
				m.currentView = viewDLTFetch
				return m, m.dltFetch.Init()
			case viewHistory:
				m.dltHistory = newDLTHistoryModel(m.dataDir)
				m.currentView = viewDLTHistory
				return m, m.dltHistory.Init()
			case viewPredict:
				m.dltPredict = newDLTPredictModel(m.dataDir)
				m.currentView = viewDLTPredict
				return m, m.dltPredict.Init()
			}
		}
		// SSQ 视图
		switch target {
		case viewFetch:
			m.fetch = newFetchModel(m.dataDir)
			m.currentView = viewFetch
			return m, m.fetch.Init()
		case viewHistory:
			m.history = newHistoryModel(m.dataDir)
			m.currentView = viewHistory
			return m, m.history.Init()
		case viewPredict:
			m.predict = newPredictModel(m.dataDir)
			m.currentView = viewPredict
			return m, m.predict.Init()
		}
	}

	var cmd tea.Cmd
	switch m.currentView {
	case viewLotterySelect:
		m.lotterySelect, cmd = m.lotterySelect.Update(msg)
	case viewMenu:
		m.menu, cmd = m.menu.Update(msg)
	case viewHistory:
		m.history, cmd = m.history.Update(msg)
	case viewPredict:
		m.predict, cmd = m.predict.Update(msg)
	case viewFetch:
		m.fetch, cmd = m.fetch.Update(msg)
	case viewDLTHistory:
		m.dltHistory, cmd = m.dltHistory.Update(msg)
	case viewDLTPredict:
		m.dltPredict, cmd = m.dltPredict.Update(msg)
	case viewDLTFetch:
		m.dltFetch, cmd = m.dltFetch.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.currentView {
	case viewLotterySelect:
		return m.lotterySelect.View()
	case viewMenu:
		return m.menu.View()
	case viewHistory:
		return m.history.View()
	case viewPredict:
		return m.predict.View()
	case viewFetch:
		return m.fetch.View()
	case viewDLTHistory:
		return m.dltHistory.View()
	case viewDLTPredict:
		return m.dltPredict.View()
	case viewDLTFetch:
		return m.dltFetch.View()
	}
	return ""
}

// Run 初始化并启动 TUI
func Run(dataDir string) error {
	p := tea.NewProgram(newModel(dataDir), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
