package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/madchin/go-cli-dirview/actual_path"
	file_traversal "github.com/madchin/go-cli-dirview/file_traversal"
	"github.com/madchin/go-cli-dirview/help"
)

type state int

const (
	traversalState state = iota
	writingState
	helpState
)

type model struct {
	viewport      viewport
	actualPath    tea.Model
	fileTraversal tea.Model
	input         tea.Model
	help          tea.Model
	state         state
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.input.Init(), m.fileTraversal.Init(), m.actualPath.Init())
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport = determineViewport(msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "/":
			if m.state == writingState || m.state == traversalState {
				m.input, cmd = m.input.Update(msg)
			}
			if m.state == writingState {
				m.state = traversalState
				return m, cmd
			}
			if m.state == traversalState || m.state == helpState {
				m.state = writingState
				return m, cmd
			}
			return m, cmd
		}
	case help.Leave:
		m.state = writingState
		return m, cmd
	case help.Display:
		m.state = helpState
		return m, cmd
	case actual_path.ActualPath:
		m.actualPath, cmd = m.actualPath.Update(msg)
		return m, cmd
	}
	if m.state == traversalState {
		m.fileTraversal, cmd = m.fileTraversal.Update(msg)
		return m, cmd
	}
	if m.state == writingState {
		m.input, cmd = m.input.Update(msg)
	}
	if m.state == helpState {
		m.help, cmd = m.help.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	var content string
	if m.state == helpState {
		content = m.help.(help.Model).WithContent(
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
		return m.viewport.renderView(content, m.help.(help.Model).Cursor)
	}
	return m.actualPath.View() + m.viewport.renderView(m.fileTraversal.View(), m.fileTraversal.(file_traversal.Model).Cursor) + "\n\nPress ctrl+c to quit.\n\n" + m.input.View()
}
