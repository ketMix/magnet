package data

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
	Kind     CellKind
	Polarity Polarity
	// ???
}

func NewLevel(path string) (Level, error) {
	level := Level{}
	err := level.LoadFromFile(path)
	return level, err
}
