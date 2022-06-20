package game

import (
	"math"
)

type PhysicsObject struct {
	// X and Y position in the world.
	X, Y float64

	// Motion Vector
	vX, vY float64

	// Magnetism
	polarity       Polarity
	magnetic       bool
	magnetStrength float64
	magnetRadius   float64

	// Size
	radius float64
}

// Retrieve the attractive/repulsive/neutral vector
func (p *PhysicsObject) GetMagneticVector(t PhysicsObject) (float64, float64) {
	// Get the magnitude of the distance to target
	vX, vY := GetDistanceVector(p.X, p.Y, t.X, t.Y)
	magnitude := GetMagnitude(vX, vY)
	vX, vY = Normalize(vX, vY, magnitude)
	affect := p.GetMagneticAffect(t.polarity, magnitude)
	return vX * affect, vY * affect
}

// Retrieve the strength of the magnetic field based on distance to target
func (p *PhysicsObject) GetMagneticAffect(polarity Polarity, magnitude float64) float64 {
	// If either source or target is neutral, return no affect
	if p.polarity == NeutralPolarity || polarity == NeutralPolarity {
		return 1.0
	}

	// Strength of magnet is inversely proportional to distance from target
	distanceRatio := p.magnetRadius / (magnitude * magnitude)
	affect := p.magnetStrength * distanceRatio

	// If they are the same polarity, produce negative (repulsive) affect
	if p.polarity != polarity {
		affect = -affect
	}

	return affect
}

// Returns distance between two points on each axis
func GetDistanceVector(sx, sy, tx, ty float64) (float64, float64) {
	return tx - sx, ty - sy
}

// Calculates the magnitude of a vector
func GetMagnitude(x, y float64) (mag float64) {
	return math.Hypot(x, y)
}

// Takes two sets of coordinates
// Returns a normalized vector for the direction from s to t
func GetDirection(sx, sy, tx, ty float64) (float64, float64) {
	vX, vY := GetDistanceVector(sx, sy, tx, ty)
	magnitude := GetMagnitude(vX, vY)
	return Normalize(vX, vY, magnitude)
}

func Normalize(vX, vY, magnitude float64) (float64, float64) {
	return (vX / magnitude), (vY / magnitude)
}

func IsWithinRadius(sx, sy, tx, ty, radius float64) bool {
	minX, maxX := sx-radius, sx+radius
	minY, maxY := sy-radius, sy+radius
	withinX := minX < tx && tx < maxX
	withinY := minY < ty && ty < maxY
	return withinX && withinY
}
