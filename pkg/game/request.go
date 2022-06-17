package game

// Request represents results from an entity's action completion.
type Request interface {
}

// UseToolRequest attempts to use the tool at a given cell.
type UseToolRequest struct {
	x, y int
	kind ToolKind // ???
}

// Belt-related requests.

// SelectToolbeltItemRequest selects a given toolbelt item
type SelectToolbeltItemRequest struct {
	kind ToolKind
}
