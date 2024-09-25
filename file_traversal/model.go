package file_traversal

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madchin/go-cli-dirview/actual_path"
	view "github.com/madchin/go-cli-dirview/file_traversal/internal"
	"github.com/muesli/termenv"
)

type Model struct {
	choices  view.Choices
	Cursor   int
	selected map[int]struct{}
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
		choices, err := view.ReadFileNamesViaLs()
		if err != nil {
			return globalErr{err}
		}
		return choices
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case view.Choices:
		m.Cursor = 0
		m.choices = msg
		return m, func() tea.Msg {
			return actual_path.Load()
		}
	case globalErr:
		m.choices = view.Choices{C: []string{msg.wrap.Error()}}
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.Cursor > 0 {
				m.Cursor--
				return m, nil
			}
			if m.Cursor == 0 {
				m.Cursor = len(m.choices.C) - 1
				return m, nil
			}
		case tea.KeyDown:
			if m.Cursor < len(m.choices.C)-1 {
				m.Cursor++
				return m, nil
			}
			if m.Cursor == len(m.choices.C)-1 {
				m.Cursor = 0
				return m, nil
			}
		case tea.KeyEnter, tea.KeyRight:
			return m, tea.Cmd(func() tea.Msg {
				destination := m.choices.C[m.Cursor]
				if !os.IsPathSeparator(destination[len(destination)-1]) {
					m.choices.C = []string{"You dont have permission to run this file, want to make it executable and exec it?", "* Yes", "* No"}
					return m.choices
				}
				if err := os.Chdir(destination); err != nil {
					return globalErr{err}
				}
				choices, err := view.ReadFileNamesViaLs()
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
				choices, err := view.ReadFileNamesViaLs()
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
	body := view.Body(output, m.choices.C, m.Cursor)
	total := len(body.EmptyMark) + len(body.FocusMark)
	for i := 0; i < len(m.choices.C); i++ {
		total += len(m.choices.C[i])
	}
	s.Grow(total)
	_, _ = body.Render(s.WriteString)

	return s.String()
}

func (m Model) changeCursorPosOnKeystrokePress(keystroke byte) Model {
	for i, choice := range m.choices.C {
		if i <= m.Cursor {
			continue
		}
		if choice[0] == keystroke {
			m.Cursor = i
			return m
		}
	}
	for i, choice := range m.choices.C[:m.Cursor] {
		if choice[0] == keystroke {
			m.Cursor = i
			return m
		}
	}
	return m
}
