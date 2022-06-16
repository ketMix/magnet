package game

// Request represents results from an entity's action completion.
type Request interface {
}

// SpawnTurretRequest attempts to spawn a turret of a type at a given cell.
type SpawnTurretRequest struct {
	x, y int
	kind int // ???
}
