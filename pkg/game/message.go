package game

// Message represents a timed message that should show on screen but disappear after a time.
type Message struct {
	lifetime  int
	deathtime int
	from      string
	content   string
}
