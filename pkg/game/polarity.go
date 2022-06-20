package game

import "image/color"

type Polarity int
type PolarityColorScale struct {
	r float64
	g float64
	b float64
	a float64
}

const (
	NegativePolarity Polarity = -1
	NeutralPolarity           = 0
	PositivePolarity          = 1
)

// Returns raw RGB values for provided polarity
func GetPolarityColor(p Polarity) color.RGBA {
	switch p {
	case NegativePolarity:
		return color.RGBA{0, 0, 255, 255}
	case PositivePolarity:
		return color.RGBA{255, 0, 0, 255}
	case NeutralPolarity:
		fallthrough
	default:
		return color.RGBA{255, 255, 255, 255}
	}
}

// Returns the scale to apply to existing color for provided polarity
func GetPolarityColorScale(p Polarity) (float64, float64, float64, float64) {
	switch p {
	case NegativePolarity:
		return .25, .25, 1, 1
	case PositivePolarity:
		return 1, .25, .25, 1
	case NeutralPolarity:
		fallthrough
	default:
		return 1.0, 1.0, 1.0, 1.0
	}
}
