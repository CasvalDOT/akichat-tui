package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/akichat/core"
	chat "github.com/akichat/core"
	"github.com/akichat/tui/models"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

type modelState struct {
	messages      []chat.Message
	lastMessageID string
}

type model struct {
	state             modelState
	sub               chan msgMessage
	chat              chat.Chat
	viewport          viewport.Model
	textarea          textarea.Model
	headerStyle       lipgloss.Style
	headerSystemStyle lipgloss.Style
	texareaStyle      lipgloss.Style
	viewportStyle     lipgloss.Style
	authorStyle       lipgloss.Style
	timeStyle         lipgloss.Style
	err               error
}

type msgMessage struct {
	messages []chat.Message
}

func initTextarea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.SetWidth(70)
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	ta.Focus()

	return ta
}

func initViewport() viewport.Model {
	vp := viewport.New(70, 30)
	vp.SetContent(`Type a message and press Enter to send.`)

	return vp
}

func initialModel() model {
	ta := initTextarea()
	vp := initViewport()
	c := chat.NewChat(core.ChatTypeHentakihabara)

	return model{
		sub: make(chan msgMessage),
		state: modelState{
			messages:      []chat.Message{},
			lastMessageID: "0",
		},
		chat:          c,
		textarea:      ta,
		viewport:      vp,
		viewportStyle: lipgloss.NewStyle().Background(lipgloss.Color("#ffffff")),
		authorStyle:   lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#efefef")),
		timeStyle:     lipgloss.NewStyle().Bold(false).Background(lipgloss.Color("#efefef")),
		texareaStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#dedede")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#dedede")).
			BorderTop(true).
			BorderRight(true).
			BorderBottom(true).
			BorderLeft(true),
		headerStyle:       lipgloss.NewStyle().Width(70).Foreground(lipgloss.Color("5")).Background(lipgloss.Color("#efefef")),
		headerSystemStyle: lipgloss.NewStyle().Align(lipgloss.Center).Width(70).Foreground(lipgloss.Color("2")).Background(lipgloss.Color("#ffffff")),
		err:               nil,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		m.listenForActivity(m.sub),
		m.waitForActivity(m.sub),
	)
}

func (m model) listenForActivity(sub chan msgMessage) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Second * 2)
			messages, _ := m.chat.ReadMessages(m.state.lastMessageID)
			if len(messages) > 0 {
				m.state.lastMessageID = messages[len(messages)-1].ID
			}

			sub <- msgMessage{
				messages: messages,
			}
		}
	}
}

func (m model) waitForActivity(sub chan msgMessage) tea.Cmd {
	return func() tea.Msg {
		return msgMessage(<-sub)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch event := msg.(type) {
	case msgMessage:
		isScrollable := false
		if len(event.messages) != len(m.state.messages) {
			isScrollable = true
		}

		m.state.messages = append(m.state.messages, event.messages...)
		m.viewport.SetContent(m.render())
		if isScrollable {
			m.viewport.GotoBottom()
		}

		return m, m.waitForActivity(m.sub)

	case tea.KeyMsg:
		switch event.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			go m.chat.WriteMessage(m.textarea.Value())
			m.viewport.SetContent(m.render())
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	case errMsg:
		m.err = event
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) render() string {
	lines := []string{}

	for _, message := range m.state.messages {
		lineAuthor := fmt.Sprintf("[%s]: ", message.Author)
		lineTime := fmt.Sprintf("%s", message.Time)
		lineContent := models.MessageText(message.Content)

		var line string
		if message.Type == core.MessageSystemType {
			line = m.headerSystemStyle.Render(lineAuthor + lineTime + "\n" + lineContent.Parse() + "\n")
		} else {
			line = m.headerStyle.Render(m.authorStyle.Render(lineAuthor)+m.timeStyle.Render(lineTime)) + "\n" + lineContent.Parse() + "\n"
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewportStyle.Render(m.viewport.View()),
		m.texareaStyle.Render(m.textarea.View()),
	) + "\n\n"
}

func NewMainView() tea.Model {
	model := initialModel()
	return &model
}
