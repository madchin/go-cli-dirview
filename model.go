package main

import (
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/madchin/go-cli-dirview/view"
	"github.com/muesli/termenv"
)

const maxFileNameLen = 256

type currDir struct{ d string }

type model struct {
	terminalHeight   int
	choices          choices
	currentDirectory currDir
	cursor           int
	selected         map[int]struct{}
}

type globalErr struct {
	wrap error
}

func (r globalErr) Error() string {
	return r.wrap.Error()
}

func (m model) Init() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		chcs, err := readFileNamesViaLs()
		if err != nil {
			return globalErr{err}
		}
		return choices{chcs}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
	case choices:
		m.cursor = 0
		m.choices = msg
		return m, tea.Cmd(func() tea.Msg {
			dir, err := exec.Command("pwd").Output()
			if err != nil {
				return globalErr{err}
			}
			return currDir{string(dir)}
		})
	case currDir:
		m.currentDirectory = msg
	case globalErr:
		m.choices = choices{[]string{msg.wrap.Error()}}
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
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
				chcs, err := readFileNamesViaLs()
				if err != nil {
					return globalErr{err}
				}
				return choices{chcs}
			})
		case "left", "esc":
			return m, tea.Cmd(func() tea.Msg {
				destination := ".."
				if err := os.Chdir(destination); err != nil {
					return globalErr{err}
				}
				chcs, err := readFileNamesViaLs()
				if err != nil {
					return globalErr{err}
				}
				return choices{chcs}
			})
		case "a":
			fallthrough
		case "b":
			fallthrough
		case "c":
			fallthrough
		case "d":
			fallthrough
		case "e":
			fallthrough
		case "f":
			fallthrough
		case "g":
			fallthrough
		case "h":
			fallthrough
		case "i":
			fallthrough
		case "j":
			fallthrough
		case "k":
			fallthrough
		case "l":
			fallthrough
		case "m":
			fallthrough
		case "n":
			fallthrough
		case "o":
			fallthrough
		case "p":
			fallthrough
		case "r":
			fallthrough
		case "s":
			fallthrough
		case "t":
			fallthrough
		case "u":
			fallthrough
		case "w":
			fallthrough
		case "x":
			fallthrough
		case "y":
			fallthrough
		case "z":
			fallthrough
		case "A":
			fallthrough
		case "B":
			fallthrough
		case "C":
			fallthrough
		case "D":
			fallthrough
		case "E":
			fallthrough
		case "F":
			fallthrough
		case "G":
			fallthrough
		case "H":
			fallthrough
		case "I":
			fallthrough
		case "J":
			fallthrough
		case "K":
			fallthrough
		case "L":
			fallthrough
		case "M":
			fallthrough
		case "N":
			fallthrough
		case "O":
			fallthrough
		case "P":
			fallthrough
		case "R":
			fallthrough
		case "S":
			fallthrough
		case "T":
			fallthrough
		case "U":
			fallthrough
		case "W":
			fallthrough
		case "X":
			fallthrough
		case "Y":
			fallthrough
		case "Z":
			return m.changeCursorPosOnKeystrokePress(msg.String()[0]), nil
		}
	}
	return m, nil
}

func (m model) View() string {
	output := termenv.NewOutput(os.Stdout)
	s := strings.Builder{}
	header := view.Header(output, m.currentDirectory.d)
	body := view.Body(output, m.choices.c, m.cursor)
	footer := view.Footer()
	viewport := view.Viewport(m.terminalHeight)
	total := len(header.Content) + len(body.EmptyMark) + len(body.FocusMark) + len(footer.Content)
	for i := 0; i < len(m.choices.c); i++ {
		total += len(m.choices.c[i])
	}
	s.Grow(total)
	_, _ = header.Render(viewport, s.WriteString)
	_, _ = body.Render(viewport, s.WriteString)
	_, _ = footer.Render(viewport, s.WriteString)

	return s.String()
}

func (m model) changeCursorPosOnKeystrokePress(keystroke byte) model {
	for i, choice := range m.choices.c {
		if i <= m.cursor {
			continue
		}
		if choice[0] == keystroke {
			m.cursor = i
			return m
		}
	}
	for i, choice := range m.choices.c[:m.cursor] {
		if choice[0] == keystroke {
			m.cursor = i
			return m
		}
	}
	return m
}
