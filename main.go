package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	file_traversal "github.com/madchin/go-cli-dirview/file_traversal_model"
)

func main() {
	p := tea.NewProgram(model{fileTraversal: file_traversal.New()})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
