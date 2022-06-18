package game

import (
	"math"
)

type PhysicsObject struct {
	// X and Y position in the world.
	X, Y float64

	// Motion Vector
	vX, vY float64

	// // Magnetism
	// magnetStrength float64
	// magnetRadius   float64

	// TODO: Other stuff, like radius/box or something.
}

// Takes two sets of coordinates
// Returns a normalized vector for the direction from s to t
func GetNormalizedDirection(sx, sy, tx, ty float64) (x, y float64) {
	return NormalizeVector((tx - sx), (ty - sy))
}

func NormalizeVector(vx, vy float64) (x, y float64) {
	magnitude := math.Sqrt(vx*vx + vy*vy)
	return (vx / magnitude), (vy / magnitude)
}
