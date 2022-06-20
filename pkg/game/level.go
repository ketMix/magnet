package game

type Level struct {
	LevelConfig
}
type CellKind int

const (
	NoneCell CellKind = iota
	PlayerCell
	PathCell
	CoreCell
	BlockedCell
	EmptyCell
	NorthSpawnCell
	SouthSpawnCell
	EnemyPositiveCell // for testing
	EnemyNegativeCell // for testing
)

type Cell struct {
	kind   CellKind
	alt    bool
	entity Entity // This is used during level -> live cells construction to store any placed turrets or similar.
	// ???
}

func NewLevel(path string) (Level, error) {
	level := Level{}
	err := level.LoadFromFile(path)
	return level, err
}
