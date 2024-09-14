package view

const (
	dimensionMin = iota
	dimensionMedium
	dimensionMax
)

type viewport struct {
	dimension int
}

func Viewport(currHeight int) viewport {
	if currHeight/8 < 2 {
		return viewport{dimensionMin}
	}
	if currHeight/8 < 4 {
		return viewport{dimensionMedium}
	}
	return viewport{dimensionMax}
}
