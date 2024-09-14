package main

import (
	tea "github.com/charmbracelet/bubbletea"
	file_traversal "github.com/madchin/go-cli-dirview/file_traversal_model"
)

type model struct {
	fileTraversal file_traversal.Model
}

func (m model) Init() tea.Cmd {
	return m.fileTraversal.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.fileTraversal.Update(msg)
}

func (m model) View() string {
	return m.fileTraversal.View()
}
