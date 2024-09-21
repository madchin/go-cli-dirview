package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	file_traversal "github.com/madchin/go-cli-dirview/file_traversal"
	"github.com/madchin/go-cli-dirview/help"
	input "github.com/madchin/go-cli-dirview/input"
)

func main() {
	p := tea.NewProgram(model{fileTraversal: file_traversal.New(), input: input.New(), help: help.New()})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
