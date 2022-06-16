package game

import (
	"bufio"
	"bytes"
	"path"
	"strings"
)

type Level struct {
	title  string
	width  int
	height int
	cells  [][]Cell
}

type CellKind int

const (
	NoneCell CellKind = iota
	PlayerCell
	PathCell
	CoreCell
	BlockedCell
	NorthSpawnCell
	SouthSpawnCell
)

type Cell struct {
	kind CellKind
	// ???
}

func newCell(r rune) (c Cell) {
	switch r {
	case 'N': // north spawn
		c.kind = NorthSpawnCell
	case 'S': // south spawn
		c.kind = SouthSpawnCell
	case 'v': // pathing node
		c.kind = PathCell
	case 'C': // core
		c.kind = CoreCell
	case '@': // player spawn
		c.kind = PlayerCell
	case '#': // unbuildable tile
		c.kind = BlockedCell
	}
	return c
}

func loadLevel(p string) (l Level, err error) {
	b, err := readFile(path.Join("levels", p+".txt"))
	if err != nil {
		return l, err
	}

	parsingHeader := true

	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		t := scanner.Text()
		if parsingHeader {
			if len(t) == 0 {
				parsingHeader = false
			} else if t[0] == 'T' {
				l.title = strings.TrimSpace(t[1:])
			}
		} else {
			if len(t) > l.width {
				l.width = len(t)
			}
			l.cells = append(l.cells, []Cell{})
			for _, r := range t {
				l.cells[l.height] = append(l.cells[l.height], newCell(r))
			}
			l.height++
		}
	}

	return l, nil
}
