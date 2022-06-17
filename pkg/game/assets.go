package game

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"path"
)

//go:embed assets/*
var assets embed.FS

func readFile(p string) ([]byte, error) {
	return assets.ReadFile(path.Join("assets", p))
}

func readImage(p string) (image.Image, error) {
	data, err := readFile(path.Join("images", p))
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

func readSound(p string) (*Sound, error) {
	data, err := readFile(path.Join("sounds", p))
	if err != nil {
		return nil, err
	}
	snd, err := NewSound(data)
	return snd, err
}
