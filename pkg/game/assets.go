package game

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"path"
	"strings"
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

func readImagesByPrefix(prefix string) ([]image.Image, error) {
	var images []image.Image
	fileList, err := assets.ReadDir(path.Join("assets", "images"))
	for _, file := range fileList {
		if !file.IsDir() && strings.HasPrefix(file.Name(), prefix) {
			image, err := readImage(file.Name())
			if err != nil {
				return images, err
			}
			images = append(images, image)
		}
	}
	return images, err
}

func getPathFiles(p string) ([]string, error) {
	var files []string
	fileList, err := assets.ReadDir(path.Join("assets", p))
	for _, file := range fileList {
		if !file.IsDir() {
			files = append(files, strings.Split(file.Name(), ".")[0])
		}
	}
	return files, err
}
