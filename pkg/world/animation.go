package world

import "github.com/hajimehoshi/ebiten/v2"

// Animation manages updating and iterating through a slice of images.
type Animation struct {
	elapsed   float64
	speed     float64
	frameTime float64
	rotation  float64
	index     int
	mirror    bool
	images    []*ebiten.Image
}

// Image returns the current image frame.
func (a *Animation) Image() *ebiten.Image {
	return a.images[a.index]
}

// Update updates the animation's current image index based upon elapsed ticks.
func (a *Animation) Update() {
	// Bail if we have no frame time.
	if a.frameTime == 0 || a.speed == 0 {
		return
	}

	// Add elapsed time and iterate the frames when needed.
	a.elapsed++
	for a.elapsed >= a.frameTime*a.speed {
		a.elapsed -= a.frameTime * a.speed
		a.Iterate()
	}
}

// Iterate steps through frames and updates the current image index.
func (a *Animation) Iterate() {
	a.index++
	if a.index >= len(a.images) {
		a.index = 0
	}
}

// Draw draws the current animation image to screen.
func (a *Animation) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	aop := &ebiten.DrawImageOptions{}

	aop.GeoM.Translate(
		-float64(a.Image().Bounds().Dx())/2,
		-float64(a.Image().Bounds().Dy())/2,
	)

	if a.mirror {
		aop.GeoM.Scale(-1, 1)
	}

	aop.GeoM.Rotate(a.rotation)
	aop.GeoM.Translate(
		float64(a.Image().Bounds().Dx())/2,
		float64(a.Image().Bounds().Dy())/2,
	)

	aop.GeoM.Concat(op.GeoM)
	aop.ColorM.Concat(op.ColorM)

	// Draw from center.
	aop.GeoM.Translate(
		-float64(a.Image().Bounds().Dx())/2,
		-float64(a.Image().Bounds().Dy())/2,
	)
	// Draw to screen.
	screen.DrawImage(a.Image(), aop)
}

func NewAnimation(images []*ebiten.Image, frameTime, speed float64) *Animation {
	return &Animation{
		images:    images,
		frameTime: frameTime,
		speed:     speed,
	}
}
