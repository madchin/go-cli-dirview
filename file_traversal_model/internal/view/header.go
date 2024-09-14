package view

import (
	"fmt"

	"github.com/muesli/termenv"
)

type header struct {
	Content string
}

func Header(output *termenv.Output, data string) header {
	return header{
		Content: output.String(fmt.Sprintf("%s\n", data)).Foreground(output.Color("5")).String(),
	}
}

func (h header) Render(viewport viewport, renderer func(content string) (n int, err error)) (int, error) {
	return renderer(h.Content)
}
