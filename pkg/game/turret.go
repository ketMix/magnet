package game

type Turret struct {
	speed float64 // speed of projecticle
	rate  float64 // projecticles per second
	tick  int     // counter for fire rate
}

// Keep track of ticks for fire rate
func (t *Turret) Tick() {
	if t.tick > 0 {
		// If we've waited long enough, reset tick country
		if float64(t.tick) >= (t.rate * 60) {
			t.tick = 0
		} else {
			t.tick += 1
		}
	}
}

// If we have a reset tick counter, we can fire and start the tick counter
func (t *Turret) CanFire() bool {
	if t.tick == 0 {
		t.tick += 1
		return true
	}
	return false
}
