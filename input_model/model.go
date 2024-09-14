package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			m.input.SetValue("")
			if m.input.Focused() {
				m.input.Blur()
				return m, nil
			}
			return m, m.input.Focus()
		}
	}
	m.input, _ = m.input.Update(msg)
	return m, nil
}

func (m Model) View() string {
	m.input.Placeholder = "/help"
	return m.input.View()
}
