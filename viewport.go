package main

import "strings"

const (
	veryLow = 3
	low     = 6
	medium  = 12
	max     = 18
)

type viewport struct {
	dimension int
}

func determineViewport(currViewHeight int) viewport {
	if currViewHeight/8 < 2 {
		return viewport{veryLow}
	}
	if currViewHeight/8 < 3 {
		return viewport{low}
	}
	if currViewHeight/8 < 4 {
		return viewport{medium}
	}
	return viewport{max}
}

func (v viewport) renderView(content string, cursor int) (view string) {
	splitted := strings.Split(content, "\n")
	stop := len(splitted)
	if cursor < 0 {
		cursor = 0
	}
	if cursor > len(splitted)-1 {
		cursor = len(splitted) - 1
	}
	if stop > cursor+v.dimension {
		stop = cursor + v.dimension
		cursor = stop - v.dimension
	}
	return strings.Join(splitted[cursor:stop], "\n")
}
