package views

import (
	"strings"

	"github.com/akichat/core"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type loginViewState struct {
	username string
	password string
}

type LoginView struct {
	chat       core.Chat
	fields     []textinput.Model
	state      loginViewState
	onSubmit   func()
	focusIndex int
}

func (m LoginView) submit() {
	m.chat.Login(m.fields[0].Value(), m.fields[1].Value())
	isAuthenticated := m.chat.IsAuthenticated()
	if isAuthenticated {
		m.onSubmit()
	}
}

func (m *LoginView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.fields))

	for i := range m.fields {
		m.fields[i], cmds[i] = m.fields[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m LoginView) Init() tea.Cmd {
	return tea.Batch()
}

func (m LoginView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch event := msg.(type) {
	case tea.KeyMsg:
		switch event.Type {
		case tea.KeyEnter:
			m.submit()
			return m, tea.Batch()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab, tea.KeyDown:
			m.focusIndex = m.focusIndex + 1
			if m.focusIndex > len(m.fields)-1 {
				m.focusIndex = 0
			}

			for i := range m.fields {
				m.fields[i].Blur()
			}

			m.fields[m.focusIndex].Focus()

			return m, tea.Batch()
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m LoginView) View() string {
	var b strings.Builder

	b.WriteString("Login or enter an username for anon access\n\n")
	for i := range m.fields {
		b.WriteString(m.fields[i].View())
		if i < len(m.fields)-1 {
			b.WriteRune('\n')
			b.WriteRune('\n')
		}
	}

	return b.String()
}

func NewLoginView(onSubmit func()) tea.Model {
	fields := []textinput.Model{}

	usernameField := textinput.New()
	usernameField.Placeholder = "Username"
	fields = append(fields, usernameField)

	passwordField := textinput.New()
	passwordField.Placeholder = "Password"
	passwordField.EchoMode = textinput.EchoPassword
	fields = append(fields, passwordField)

	m := LoginView{
		fields:   fields,
		onSubmit: onSubmit,
		chat:     core.NewChat(core.ChatTypeHentakihabara),
	}

	return &m
}
