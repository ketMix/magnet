package data

import (
	"bufio"
	"bytes"
	"path"
	"strings"
)

type LevelConfig struct {
	Title   string
	Tileset string
	Width   int
	Height  int
	Cells   [][]Cell
}

func (l *LevelConfig) newCell(r rune) (c Cell) {
	switch r {
	case 'N': // north spawn
		c.Kind = NorthSpawnCell
		c.Alt = true
	case 'S': // south spawn
		c.Kind = SouthSpawnCell
	case 'v': // pathing node
		c.Kind = PathCell
	case 'C': // core
		c.Kind = CoreCell
	case '@': // player spawn
		c.Kind = PlayerCell
	case '#': // unbuildable tile
		c.Kind = BlockedCell
	case ',': // alternative open cell
		c.Alt = true
	case ' ': // unpathable tile -- like '#', but no image.
		c.Kind = EmptyCell
	case '+': // TEST - positive enemy cell
		c.Kind = EnemyPositiveCell
	case '-': // TEST - negative enemy cell
		c.Kind = EnemyNegativeCell
	}

	return c
}

func (l *LevelConfig) LoadFromFile(p string) (err error) {
	b, err := ReadFile(path.Join("levels", p+".txt"))
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
				l.Title = strings.TrimSpace(t[1:])
			} else if t[0] == 'S' {
				l.Tileset = strings.TrimSpace(t[1:])
			}
		} else {
			if len(t) > l.Width {
				l.Width = len(t)
			}
			l.Cells = append(l.Cells, []Cell{})
			for _, r := range t {
				l.Cells[l.Height] = append(l.Cells[l.Height], l.newCell(r))
			}
			l.Height++
		}
	}

	return nil
}
