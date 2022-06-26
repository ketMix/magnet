package data

import (
	"bufio"
	"bytes"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type EntityConfig struct {
	Title            string
	Points           int
	Health           int
	Speed            float64
	Radius           float64
	Damage           int
	AttackRange      float64
	AttackRate       float64
	ProjecticleSpeed float64
	Polarity         Polarity
	Magnetic         bool
	MagnetStrength   float64
	MagnetRadius     float64
	Images           []*ebiten.Image
	WalkImages       []*ebiten.Image
	HeadImages       []*ebiten.Image
}

func (e *EntityConfig) LoadFromFile(p string) error {
	b, err := ReadFile(path.Join("entities", p+".txt"))
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		t := scanner.Text()
		value := strings.TrimSpace(t[1:])
		switch t[0] {
		case 'T':
			e.Title = value
		case 'C':
			e.Points, err = strconv.Atoi(value)
		case 'H':
			e.Health, err = strconv.Atoi(value)
		case 'D':
			e.Damage, err = strconv.Atoi(value)
		case 'R':
			e.AttackRange, err = strconv.ParseFloat(value, 64)
		case 'X':
			e.AttackRate, err = strconv.ParseFloat(value, 64)
		case 'O':
			e.ProjecticleSpeed, err = strconv.ParseFloat(value, 64)
		case 'S':
			e.Speed, err = strconv.ParseFloat(value, 64)
		case 'r':
			e.Radius, err = strconv.ParseFloat(value, 64)
		case 'P':
			switch value {
			case "positive":
				e.Polarity = PositivePolarity
			case "negative":
				e.Polarity = NegativePolarity
			default:
				fallthrough
			case "neutral":
				e.Polarity = NeutralPolarity
			}
		case 'M':
			e.Magnetic, err = strconv.ParseBool(value)
		case 'Y':
			e.MagnetStrength, err = strconv.ParseFloat(value, 64)
		case 'Z':
			e.MagnetRadius, err = strconv.ParseFloat(value, 64)
		case 'I':
			// Load images using prefix in value
			images, err := ReadImagesByPrefix(value)
			if err != nil {
				return err
			}
			for _, image := range images {
				img := ebiten.NewImageFromImage(image)
				e.Images = append(e.Images, img)
			}
		case 'W':
			// Load images using prefix in value
			images, err := ReadImagesByPrefix(value)
			if err != nil {
				return err
			}
			for _, image := range images {
				img := ebiten.NewImageFromImage(image)
				e.WalkImages = append(e.WalkImages, img)
			}
		case 'i':
			// Load images using prefix in value
			images, err := ReadImagesByPrefix(value)
			if err != nil {
				return err
			}
			for _, image := range images {
				img := ebiten.NewImageFromImage(image)
				e.HeadImages = append(e.HeadImages, img)
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

func NewPlayerConfig(i int) (EntityConfig, error) {
	config := EntityConfig{}
	err := config.LoadFromFile(fmt.Sprintf("player%d", i))
	return config, err
}

func NewCoreConfig() (EntityConfig, error) {
	config := EntityConfig{}
	err := config.LoadFromFile("core")
	return config, err
}
