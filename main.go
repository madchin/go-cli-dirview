package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	file_traversal "github.com/madchin/go-cli-dirview/file_traversal_model"
	input "github.com/madchin/go-cli-dirview/input_model"
)

func main() {
	p := tea.NewProgram(model{fileTraversal: file_traversal.New(), input: input.New()})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
