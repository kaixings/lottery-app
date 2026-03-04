package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"lottery-app/internal/fetcher"
	"lottery-app/internal/storage"
)

type dltFetchModel struct {
	state   fetchState
	message string
	dataDir string
}

type dltFetchDoneMsg struct {
	count int
	err   error
}

func newDLTFetchModel(dataDir string) dltFetchModel {
	return dltFetchModel{state: fetchIdle, dataDir: dataDir}
}

func (m dltFetchModel) doFetchCmd() tea.Cmd {
	return func() tea.Msg {
		records, err := fetcher.FetchDLT(100)
		if err != nil {
			return dltFetchDoneMsg{err: err}
		}
		store, err := storage.NewDLTStore(m.dataDir)
		if err != nil {
			return dltFetchDoneMsg{err: err}
		}
		existing, _ := store.Load()
		merged := store.Merge(existing, records)
		if err := store.Save(merged); err != nil {
			return dltFetchDoneMsg{err: err}
		}
		return dltFetchDoneMsg{count: len(records)}
	}
}

func (m dltFetchModel) Init() tea.Cmd { return nil }

func (m dltFetchModel) Update(msg tea.Msg) (dltFetchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case dltFetchDoneMsg:
		if msg.err != nil {
			m.state = fetchError
			m.message = msg.err.Error()
		} else {
			m.state = fetchDone
			m.message = ""
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.state == fetchIdle || m.state == fetchDone || m.state == fetchError {
				m.state = fetchFetching
				m.message = ""
				return m, m.doFetchCmd()
			}
		case "esc", "q":
			if m.state != fetchFetching {
				m.state = fetchIdle
				return m, func() tea.Msg { return switchViewMsg{viewMenu} }
			}
		}
	}
	return m, nil
}

func (m dltFetchModel) View() string {
	s := titleStyle.Render("抓取数据（大乐透）") + "\n\n"
	switch m.state {
	case fetchIdle:
		s += normalStyle.Render("按 Enter 开始从500彩票网抓取最新大乐透数据") + "\n"
	case fetchFetching:
		s += dimStyle.Render("正在抓取数据，请稍候...") + "\n"
	case fetchDone:
		s += successStyle.Render("抓取成功！数据已保存。") + "\n"
	case fetchError:
		s += errorStyle.Render("抓取失败: "+m.message) + "\n"
	}
	s += "\n" + dimStyle.Render("Enter 开始抓取  Esc/q 返回")
	return s
}
