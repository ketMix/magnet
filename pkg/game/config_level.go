package game

import (
	"bufio"
	"bytes"
	"path"
	"strings"
)

type LevelConfig struct {
	title   string
	tileset string
	width   int
	height  int
	cells   [][]Cell
}

func (l *LevelConfig) newCell(r rune) (c Cell) {
	switch r {
	case 'N': // north spawn
		c.kind = NorthSpawnCell
		c.alt = true
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
	case ',': // alternative open cell
		c.alt = true
	case ' ': // unpathable tile -- like '#', but no image.
		c.kind = EmptyCell
	case '+': // TEST - positive enemy cell
		c.kind = EnemyPositiveCell
	case '-': // TEST - negative enemy cell
		c.kind = EnemyNegativeCell
	}

	return c
}

func (l *LevelConfig) LoadFromFile(p string) (err error) {
	b, err := readFile(path.Join("levels", p+".txt"))
	if err != nil {
		return err
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
			} else if t[0] == 'S' {
				l.tileset = strings.TrimSpace(t[1:])
			}
		} else {
			if len(t) > l.width {
				l.width = len(t)
			}
			l.cells = append(l.cells, []Cell{})
			for _, r := range t {
				l.cells[l.height] = append(l.cells[l.height], l.newCell(r))
			}
			l.height++
		}
	}

	return nil
}
