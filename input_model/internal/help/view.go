package help

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Leave struct{}

type model struct {
	content []string
}

func New(content ...string) model {
	return model{content}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, func() tea.Msg {
			return Leave{}
		}
	}
	return m, nil
}

func (m model) View() string {
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
