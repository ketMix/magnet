package game

import (
	"bufio"
	"bytes"
	"path"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type EntityConfig struct {
	title            string
	health           int
	speed            float64
	radius           float64
	damage           int
	attackRange      float64
	attackRate       float64
	projecticleSpeed float64
	polarity         Polarity
	magnetic         bool
	magnetStrength   float64
	magnetRadius     float64
	images           []*ebiten.Image
	headImages       []*ebiten.Image
}

func (e *EntityConfig) LoadFromFile(p string) error {
	b, err := readFile(path.Join("entities", p+".txt"))
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		t := scanner.Text()
		value := strings.TrimSpace(t[1:])
		switch t[0] {
		case 'T':
			e.title = value
		case 'H':
			e.health, err = strconv.Atoi(value)
		case 'D':
			e.damage, err = strconv.Atoi(value)
		case 'R':
			e.attackRange, err = strconv.ParseFloat(value, 64)
		case 'X':
			e.attackRate, err = strconv.ParseFloat(value, 64)
		case 'O':
			e.projecticleSpeed, err = strconv.ParseFloat(value, 64)
		case 'S':
			e.speed, err = strconv.ParseFloat(value, 64)
		case 'r':
			e.radius, err = strconv.ParseFloat(value, 64)
		case 'P':
			switch value {
			case "positive":
				e.polarity = PositivePolarity
			case "negative":
				e.polarity = NegativePolarity
			default:
				fallthrough
			case "neutral":
				e.polarity = NeutralPolarity
			}
		case 'M':
			e.magnetic, err = strconv.ParseBool(value)
		case 'Y':
			e.magnetStrength, err = strconv.ParseFloat(value, 64)
		case 'Z':
			e.magnetRadius, err = strconv.ParseFloat(value, 64)
		case 'I':
			// Load images using prefix in value
			images, err := readImagesByPrefix(value)
			if err != nil {
				return err
			}
			for _, image := range images {
				img := ebiten.NewImageFromImage(image)
				e.images = append(e.images, img)
			}
		case 'i':
			// Load images using prefix in value
			images, err := readImagesByPrefix(value)
			if err != nil {
				return err
			}
			for _, image := range images {
				img := ebiten.NewImageFromImage(image)
				e.headImages = append(e.headImages, img)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func NewEnemyConfig(p string) (EntityConfig, error) {
	config := EntityConfig{}
	err := config.LoadFromFile(path.Join("enemies", p))
	return config, err
}

func NewTurretConfig(p string) (EntityConfig, error) {
	config := EntityConfig{}
	err := config.LoadFromFile(path.Join("turrets", p))
	return config, err
}

func NewPlayerConfig() (EntityConfig, error) {
	config := EntityConfig{}
	err := config.LoadFromFile("player")
	return config, err
}
