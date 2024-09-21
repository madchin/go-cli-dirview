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

func (b body) Render(viewport viewport, renderer func(content string) (n int, err error)) (n int, err error) {
	start, stop := 0, len(b.data)
	if viewport.dimension == dimensionVeryLow {
		if stop > 3 {
			start, stop = b.viewRange(3)
		}
	} else if viewport.dimension == dimensionLow {
		if stop > 6 {
			start, stop = b.viewRange(6)
		}
	} else if viewport.dimension == dimensionMedium {
		if stop > 12 {
			start, stop = b.viewRange(12)
		}
	} else if viewport.dimension == dimensionMax {
		if stop > 18 {
			start, stop = b.viewRange(18)
		}
	}
	for i, line := range b.data[start:stop] {
		cursor := b.EmptyMark
		if b.focusElement == i+start {
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

func (b body) viewRange(threshold int) (start int, stop int) {
	start = b.focusElement
	if len(b.data) > start+threshold {
		stop = start + threshold
	} else {
		stop = len(b.data)
	}
	start = stop - threshold
	return start, stop
}

func isDirectory(candidate string) bool {
	return os.IsPathSeparator(candidate[len(candidate)-1])
}
