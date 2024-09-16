package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madchin/go-cli-dirview/input_model/internal/help"
)

type state int

const (
	baseState state = iota
	helpState
)

type Model struct {
	state state
	input textinput.Model
	help  tea.Model
}

func New() Model {
	return Model{baseState, textinput.New(), help.New("/help -- print usage")}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case help.Leave:
		m.state = baseState
		return m, nil
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
				m.state = helpState
				return m, cmd
			}
		}
	}
	if m.state == helpState {
		m.help, cmd = m.help.Update(msg)
	}
	if m.state == baseState {
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	m.input.Placeholder = "/help"
	if m.state == baseState {
		return m.input.View()
	}
	if m.state == helpState {
		return m.help.View()
	}
	return ""
}
