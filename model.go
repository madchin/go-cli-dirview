package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madchin/go-cli-dirview/help"
)

type state int

const (
	traversalState state = iota
	writingState
	helpState
)

type model struct {
	fileTraversal tea.Model
	input         tea.Model
	help          tea.Model
	state         state
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.input.Init(), m.fileTraversal.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "/":
			if m.state == writingState {
				m.state = m.setState(traversalState)
			} else {
				m.state = writingState
			}
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	case help.Leave:
		m.state = m.setState(writingState)
		return m, cmd
	case help.Display:
		m.state = m.setState(helpState)
		return m, cmd
	}
	if m.state == traversalState {
		m.fileTraversal, cmd = m.fileTraversal.Update(msg)
	}
	if m.state == writingState {
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.state == helpState {
		return m.help.(help.Model).WithContent(
			`Welcome to /help page!
			Overall, using dirview, you can switch between three modes, either interactive and not.
			Writing, Traversal and Help ones.
			
			Traversal mode is used to view files / directories in file tree.
				You can type letter to jump onto files/directories at this letter.
				Jump between directories using arrows (< -- jump back) (> -- jump forward)

			Writing mode is used to write commands onto application
				Actual write commands: 
					/help
			
			Help mode which displays this page.

			To toggle between traversal / writing modes press / keystroke.
			To jump onto help mode, type /help in writing mode
		`).View()
	} else {
		return fmt.Sprintf("%s\nPress ctrl+c to quit.\n%s\n", m.fileTraversal.View(), m.input.View())
	}
}

func (m model) setState(state state) state {
	m.state = state
	return m.state
}
