package help

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Leave struct{}
type Display struct{}

type Model struct {
	content []string
}

func New() Model {
	return Model{}
}

func (m Model) WithContent(content ...string) Model {
	m.content = content
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, func() tea.Msg {
		return Leave{}
	}
}

func (m Model) View() string {
	var strBuilder strings.Builder
	var total int
	for i := 0; i < len(m.content); i++ {
		total += len(m.content[i])
	}
	strBuilder.Grow(total)
	for i := 0; i < len(m.content); i++ {
		_, _ = strBuilder.WriteString(m.content[i])
	}
	return strBuilder.String()
}
