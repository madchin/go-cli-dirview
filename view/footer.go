package view

type footer struct {
	Content string
}

func Footer() footer {
	return footer{Content: "\nPress ctrl+c to quit.\n"}
}

func (f footer) Render(viewport viewport, renderer func(content string) (n int, err error)) (int, error) {
	return renderer(f.Content)
}
