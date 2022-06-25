package data

import (
	"bufio"
	"bytes"
	"path"
	"strconv"
	"strings"
)

type LevelConfig struct {
	Title   string
	Next    string
	Tileset string
	Width   int
	Height  int
	Cells   [][]Cell
	Waves   []*Wave
}

func (l *LevelConfig) newCell(r rune) (c Cell) {
	switch r {
	case 'N': // north spawn
		c.Kind = NorthSpawnCell
		c.Polarity = NegativePolarity
	case 'S': // south spawn
		c.Kind = SouthSpawnCell
		c.Polarity = PositivePolarity
	case 'v': // pathing node
		c.Kind = PathCell
	case 'C': // core
		c.Kind = CoreCell
	case '@': // player spawn
		c.Kind = PlayerCell
	case '#': // unbuildable tile
		c.Kind = BlockedCell
	case '.': // open negative cell
		c.Polarity = NegativePolarity
	case ',': // open positive cell
		c.Polarity = PositivePolarity
	case '_': // open neutral cell
		c.Polarity = NeutralPolarity
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
			} else if t[0] == 'N' {
				l.Next = strings.TrimSpace(t[1:])
			} else if t[0] == 'W' {
				s := strings.TrimSpace(t[1:])
				waveStrs := strings.Split(s, ";")
				var firstWave *Wave
				var lastWave *Wave
				for _, w := range waveStrs {
					var lastSpawn *SpawnList
					wave := &Wave{}
					spawnStrs := strings.Split(w, ",")
					for _, s := range spawnStrs {
						var amount int
						var tickDelay int
						var enemies []string
						amountAndList := strings.Split(s, " ")
						amountStr := amountAndList[0]
						// Get amount and tick delay.
						amountAndTickDelayStrs := strings.Split(amountStr, "@")
						if len(amountAndTickDelayStrs) == 1 {
							// No tick specified.
							amount, _ = strconv.Atoi(amountAndTickDelayStrs[0])
							tickDelay = 20
						} else {
							amount, _ = strconv.Atoi(amountAndTickDelayStrs[0])
							tickDelay, _ = strconv.Atoi(amountAndTickDelayStrs[1])
						}
						// Get enemies list.
						listStr := amountAndList[1]
						enemies = strings.Split(listStr, "&")
						sl := &SpawnList{
							Kinds:     enemies,
							Spawnrate: tickDelay,
							Count:     amount,
						}
						if lastSpawn == nil {
							wave.Spawns = sl
						} else {
							lastSpawn.Next = sl
						}
						lastSpawn = sl
					}
					if firstWave == nil {
						firstWave = wave
					}
					if lastWave != nil {
						lastWave.Next = wave
					}
					lastWave = wave
				}
				l.Waves = append(l.Waves, firstWave)
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
