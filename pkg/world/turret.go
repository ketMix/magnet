package world

type Turret struct {
	damage         int     // the damage of each projecticle
	projecticleNum int     // the number of projecticles per attack
	speed          float64 // speed of projecticle
	rate           float64 // projecticles per second
	tick           float64 // counter for fire rate
	attackRange    float64 // the range of the turret
}

// Keep track of ticks for fire rate
func (t *Turret) Tick(speed float64) {
	if t.tick > 0 {
		// If we've waited long enough, reset tick country
		if float64(t.tick) >= (t.rate * 60) {
			t.tick = 0
		} else {
			t.tick += speed
		}
	}
}

// If we have a reset tick counter, we can fire and start the tick counter
func (t *Turret) CanFire(speed float64) bool {
	if t.tick == 0 {
		t.tick += speed
		return true
	}
	return false
}
