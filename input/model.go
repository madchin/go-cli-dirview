package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madchin/go-cli-dirview/help"
)

type Model struct {
	input textinput.Model
}

func New() Model {
	return Model{textinput.New()}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			if m.input.Focused() {
				m.input.Blur()
				m.input.SetValue("")
				return m, nil
			}
			m.input.SetValue("/")
			return m, m.input.Focus()
		case "enter":
			switch m.input.Value() {
			case "/help":
				return m, func() tea.Msg {
					return help.Display{}
				}
			}
		default:
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	}
	return m, cmd
}

func (m Model) View() string {
	m.input.Placeholder = "Press / to start writing. For more help type /help"
	return m.input.View()
}
