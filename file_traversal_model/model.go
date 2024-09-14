package file_traversal

import (
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	data "github.com/madchin/go-cli-dirview/file_traversal_model/internal"
	"github.com/madchin/go-cli-dirview/file_traversal_model/internal/view"
	"github.com/muesli/termenv"
)

type currentDirectory struct{ d string }

type Model struct {
	terminalHeight   int
	choices          data.Choices
	currentDirectory currentDirectory
	cursor           int
	selected         map[int]struct{}
}

func New() Model {
	return Model{selected: make(map[int]struct{}, 1)}
}

type globalErr struct {
	wrap error
}

func (r globalErr) Error() string {
	return r.wrap.Error()
}

func (m Model) Init() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		choices, err := data.ReadFileNamesViaLs()
		if err != nil {
			return globalErr{err}
		}
		return choices
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
	case data.Choices:
		m.cursor = 0
		m.choices = msg
		return m, tea.Cmd(func() tea.Msg {
			dir, err := exec.Command("pwd").Output()
			if err != nil {
				return globalErr{err}
			}
			return currentDirectory{string(dir)}
		})
	case currentDirectory:
		m.currentDirectory = msg
	case globalErr:
		m.choices = data.Choices{C: []string{msg.wrap.Error()}}
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
				return m, nil
			}
			if m.cursor == 0 {
				m.cursor = len(m.choices.C) - 1
				return m, nil
			}
		case tea.KeyDown:
			if m.cursor < len(m.choices.C)-1 {
				m.cursor++
				return m, nil
			}
			if m.cursor == len(m.choices.C)-1 {
				m.cursor = 0
				return m, nil
			}
		case tea.KeyEnter, tea.KeyRight:
			return m, tea.Cmd(func() tea.Msg {
				destination := m.choices.C[m.cursor]
				if !os.IsPathSeparator(destination[len(destination)-1]) {
					m.choices.C = []string{"You dont have permission to run this file, want to make it executable and exec it?", "* Yes", "* No"}
					return m.choices
				}
				if err := os.Chdir(destination); err != nil {
					return globalErr{err}
				}
				choices, err := data.ReadFileNamesViaLs()
				if err != nil {
					return globalErr{err}
				}
				return choices
			})
		case tea.KeyEsc, tea.KeyLeft:
			return m, tea.Cmd(func() tea.Msg {
				destination := ".."
				if err := os.Chdir(destination); err != nil {
					return globalErr{err}
				}
				choices, err := data.ReadFileNamesViaLs()
				if err != nil {
					return globalErr{err}
				}
				return choices
			})
		}
		switch msg.String() {
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "r", "s", "t", "u", "w", "x", "y", "z":
			fallthrough
		case "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "R", "S", "T", "U", "W", "X", "Y", "Z":
			return m.changeCursorPosOnKeystrokePress(msg.String()[0]), nil
		}
	}
	return m, nil
}

func (m Model) View() string {
	output := termenv.NewOutput(os.Stdout)
	s := strings.Builder{}
	header := view.Header(output, m.currentDirectory.d)
	body := view.Body(output, m.choices.C, m.cursor)
	viewport := view.Viewport(m.terminalHeight)
	total := len(header.Content) + len(body.EmptyMark) + len(body.FocusMark)
	for i := 0; i < len(m.choices.C); i++ {
		total += len(m.choices.C[i])
	}
	s.Grow(total)
	_, _ = header.Render(viewport, s.WriteString)
	_, _ = body.Render(viewport, s.WriteString)

	return s.String()
}

func (m Model) changeCursorPosOnKeystrokePress(keystroke byte) Model {
	for i, choice := range m.choices.C {
		if i <= m.cursor {
			continue
		}
		if choice[0] == keystroke {
			m.cursor = i
			return m
		}
	}
	for i, choice := range m.choices.C[:m.cursor] {
		if choice[0] == keystroke {
			m.cursor = i
			return m
		}
	}
	return m
}
