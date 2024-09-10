package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

const maxFileNameLen = 256

type model struct {
	terminalHeight int
	choices        choices
	cursor         int
	selected       map[int]struct{}
}

type globalErr struct {
	wrap error
}

func (r globalErr) Error() string {
	return r.wrap.Error()
}

func (m model) Init() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		choices, err := readFileNamesViaLs()
		if err != nil {
			return globalErr{err}
		}
		return choices
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
	case choices:
		m.cursor = 0
		m.choices = msg
	case globalErr:
		m.choices = choices{[]string{msg.wrap.Error()}}
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
				return m, nil
			}
			if m.cursor == 0 {
				m.cursor = len(m.choices.c) - 1
				return m, nil
			}
		case "down":
			if m.cursor < len(m.choices.c)-1 {
				m.cursor++
				return m, nil
			}
			if m.cursor == len(m.choices.c)-1 {
				m.cursor = 0
				return m, nil
			}
		case "enter", "right":
			return m, tea.Cmd(func() tea.Msg {
				destination := m.choices.c[m.cursor]
				if !os.IsPathSeparator(destination[len(destination)-1]) {
					m.choices.c = []string{"You dont have permission to run this file, want to make it executable and exec it?", "* Yes", "* No"}
					return m.choices
				}
				if err := os.Chdir(destination); err != nil {
					return globalErr{err}
				}
				choices, err := readFileNamesViaLs()
				if err != nil {
					return globalErr{err}
				}
				m.choices = choices
				return m.choices
			})
		case "left", "esc":
			return m, tea.Cmd(func() tea.Msg {
				destination := ".."
				if err := os.Chdir(destination); err != nil {
					return globalErr{err}
				}
				choices, err := readFileNamesViaLs()
				if err != nil {
					return globalErr{err}
				}
				m.choices = choices
				return m.choices
			})
		case "t":
			return m, tea.Cmd(func() tea.Msg {
				str := os.Environ()
				fmt.Printf("str :%v", str)
				return struct{}{}
			})
		}
	}

	return m, nil
}

func (m model) View() string {
	output := termenv.NewOutput(os.Stdout)
	s := strings.Builder{}
	header := "Where do you want to head?\n\n"
	footer := "\nPress q to quit.\n"
	emptyMark := " "
	focusMark := ">"
	total := len(header) + len(emptyMark) + len(focusMark) + len(footer)
	for i := 0; i < len(m.choices.c); i++ {
		total += len(m.choices.c[i])
	}
	s.Grow(total)
	_, _ = s.WriteString(header)
	start, stop := 0, len(m.choices.c)
	if stop > m.terminalHeight/2 {
		if m.cursor == 0 {
			start = m.cursor
		} else {
			start = m.cursor - 1
		}
		if start+m.terminalHeight/2 > stop {
			start = stop - m.terminalHeight/2
		}
		stop = start + m.terminalHeight/2
	}
	for i, choice := range m.choices.c[start:stop] {
		cursor := emptyMark
		if m.cursor == i+start {
			cursor = output.String(focusMark).Foreground(output.Color("118")).String()
		}
		if os.IsPathSeparator(choice[len(choice)-1]) {
			choice = output.String(choice).Foreground(output.Color("214")).String()
		} else {
			choice = output.String(choice).Foreground(output.Color("33")).String()
		}
		_, _ = s.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}
	_, _ = s.WriteString(footer)

	return s.String()
}
