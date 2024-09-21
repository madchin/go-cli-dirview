package view

const (
	dimensionVeryLow = iota
	dimensionLow
	dimensionMedium
	dimensionMax
)

type viewport struct {
	dimension int
}

func Viewport(currHeight int) viewport {
	if currHeight/8 < 2 {
		return viewport{dimensionVeryLow}
	}
	if currHeight/8 < 3 {
		return viewport{dimensionLow}
	}
	if currHeight/8 < 4 {
		return viewport{dimensionMedium}
	}
	return viewport{dimensionMax}
}
