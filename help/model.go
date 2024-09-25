package help

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Leave struct{}
type Display struct{}

type Model struct {
	content string
	Cursor  int
}

func New() Model {
	return Model{}
}

func (m Model) WithContent(content string) Model {
	m.content = content
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyDown {
			m.Cursor++
			return m, nil
		}
		if msg.Type == tea.KeyUp {
			m.Cursor--
			return m, nil
		}
		return m, func() tea.Msg {
			return Leave{}
		}
	}
	return m, nil
}

func (m Model) View() string {
	var strBuilder strings.Builder
	var total int
	for i := 0; i < len(m.content); i++ {
		total += len(m.content)
	}
	strBuilder.Grow(total)
	for i := 0; i < len(m.content); i++ {
		_, _ = strBuilder.WriteString(m.content)
	}
	return strBuilder.String()
}
