package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"lottery-app/internal/fetcher"
	"lottery-app/internal/storage"
)

type fetchState int

const (
	fetchIdle fetchState = iota
	fetchFetching
	fetchDone
	fetchError
)

type fetchModel struct {
	state   fetchState
	message string
	dataDir string
}

type fetchDoneMsg struct {
	count int
	err   error
}

func newFetchModel(dataDir string) fetchModel {
	return fetchModel{state: fetchIdle, dataDir: dataDir}
}

func (m fetchModel) doFetchCmd() tea.Cmd {
	return func() tea.Msg {
		records, err := fetcher.Fetch(100)
		if err != nil {
			return fetchDoneMsg{err: err}
		}
		store, err := storage.NewStore(m.dataDir)
		if err != nil {
			return fetchDoneMsg{err: err}
		}
		existing, _ := store.Load()
		merged := store.Merge(existing, records)
		if err := store.Save(merged); err != nil {
			return fetchDoneMsg{err: err}
		}
		return fetchDoneMsg{count: len(records)}
	}
}

func (m fetchModel) Init() tea.Cmd { return nil }

func (m fetchModel) Update(msg tea.Msg) (fetchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case fetchDoneMsg:
		if msg.err != nil {
			m.state = fetchError
			m.message = msg.err.Error()
		} else {
			m.state = fetchDone
			m.message = ""
			_ = msg.count
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

func (m fetchModel) View() string {
	s := titleStyle.Render("抓取数据") + "\n\n"
	switch m.state {
	case fetchIdle:
		s += normalStyle.Render("按 Enter 开始从500彩票网抓取最新数据") + "\n"
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
