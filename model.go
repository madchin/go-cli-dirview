package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	traversal state = iota
	writing
)

type model struct {
	fileTraversal tea.Model
	input         tea.Model
	state         state
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.input.Init(), m.fileTraversal.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		switch msg.String() {
		case "/":
			m.state = m.toggleState()
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	}
	if m.state == traversal {
		m.fileTraversal, cmd = m.fileTraversal.Update(msg)
	}
	if m.state == writing {
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("%s\nPress ctrl+c to quit.\n%s\n", m.fileTraversal.View(), m.input.View())
}

func (m model) toggleState() state {
	if m.state == traversal {
		return writing
	}
	return traversal
}
