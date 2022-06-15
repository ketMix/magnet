package game

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"path/filepath"
)

//go:embed assets/*
var assets embed.FS

func readFile(p string) ([]byte, error) {
	return assets.ReadFile(filepath.Join("assets", p))
}

func readImage(p string) (image.Image, error) {
	data, err := readFile(filepath.Join("images", p))
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}
