package view

import (
	"fmt"
	"os"

	"github.com/muesli/termenv"
)

type body struct {
	output       *termenv.Output
	data         []string
	focusElement int
	FocusMark    string
	EmptyMark    string
	Size         int
}

func Body(
	output *termenv.Output,
	data []string,
	chosen int,
) body {
	return body{
		FocusMark:    ">",
		EmptyMark:    " ",
		output:       output,
		data:         data,
		focusElement: chosen,
	}
}

func (b body) Render(renderer func(content string) (n int, err error)) (n int, err error) {
	for i, line := range b.data {
		cursor := b.EmptyMark
		if b.focusElement == i {
			cursor = b.output.String(b.FocusMark).Foreground(b.output.Color("118")).String()
		}
		if isDirectory(line) {
			line = b.output.String(line).Foreground(b.output.Color("214")).String()
		} else {
			line = b.output.String(line).Foreground(b.output.Color("33")).String()
		}
		count, err := renderer(fmt.Sprintf("%s %s\n", cursor, line))
		if err != nil {
			return -1, err
		}
		n += count
	}
	return n, err
}

func isDirectory(candidate string) bool {
	return os.IsPathSeparator(candidate[len(candidate)-1])
}
