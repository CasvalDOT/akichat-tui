package main

import (
	"log"

	core "github.com/CasvalDOT/akichat-core"
	"github.com/CasvalDOT/akichat-tui/views"
	tea "github.com/charmbracelet/bubbletea"
)

type root struct {
	sub  chan struct{}
	chat core.Chat
	view tea.Model
}

type (
	onLogin     struct{}
	onLogout    struct{}
	activityMsg struct{}
)

func initialModel() root {
	r := root{
		sub:  make(chan struct{}),
		chat: core.NewChat(core.ChatTypeHentakihabara),
	}

	if r.chat.IsAuthenticated() == false {
		r.initUnauthenticatedView()
	} else {
		r.initAuthenticatedView()
	}

	return r
}

func (m root) waitForLogin(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return onLogin(<-sub)
	}
}

func (m *root) initAuthenticatedView() {
	mainView := views.NewMainView()

	m.view = mainView
}

func (m *root) initUnauthenticatedView() {
	loginView := views.NewLoginView(
		func() {
			m.sub <- onLogin{}
		},
	)

	m.view = loginView
}

func (m root) Init() tea.Cmd {
	return tea.Batch(
		m.view.Init(),
		m.waitForLogin(m.sub))
}

func (m root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch event := msg.(type) {
	case onLogin:
		m.initAuthenticatedView()
		return m, tea.Batch(m.view.Init(), m.waitForLogin(m.sub))
	case onLogout:
		return m, tea.Batch(m.view.Init())
	case tea.KeyMsg:
		switch event.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			cmds = append(cmds, tea.Quit)
		}
	}

	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m root) View() string {
	return m.view.View()
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
