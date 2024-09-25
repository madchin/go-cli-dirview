package actual_path

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type ActualPath struct{ path []byte }

type Model struct {
	path []byte
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return Load()
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(ActualPath); ok {
		m.path = msg.path
	}
	return m, nil
}

func (m Model) View() string {
	output := termenv.NewOutput(os.Stdout)
	return output.String(fmt.Sprintf("%s\n", m.path)).Foreground(output.Color("5")).String()
}

func Load() ActualPath {
	dir, err := exec.Command("pwd").Output()
	if err != nil {
		return ActualPath{}
	}
	return ActualPath{dir}
}
