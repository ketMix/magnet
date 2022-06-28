package data

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

func ReadFile(p string) ([]byte, error) {
	return assets.ReadFile(path.Join("assets", p))
}

func ReadImage(p string) (image.Image, error) {
	data, err := ReadFile(path.Join("images", p))
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

func ReadSound(p string) (*Sound, error) {
	data, err := ReadFile(path.Join("sounds", p))
	if err != nil {
		return nil, err
	}
	snd, err := NewSound(data)
	return snd, err
}

func ReadMusic(p string) (*Sound, error) {
	data, err := ReadFile(path.Join("music", p))
	if err != nil {
		return nil, err
	}
	bgm, err := NewSound(data)
	return bgm, err
}

func ReadImagesByPrefix(prefix string) ([]image.Image, error) {
	var images []image.Image
	fileList, err := assets.ReadDir(path.Join("assets", "images"))
	for _, file := range fileList {
		if !file.IsDir() && strings.HasPrefix(file.Name(), prefix) {
			image, err := ReadImage(file.Name())
			if err != nil {
				return images, err
			}
			images = append(images, image)
		}
	}
	return images, err
}

func GetPathFiles(p string) ([]string, error) {
	var files []string
	fileList, err := assets.ReadDir(path.Join("assets", p))
	for _, file := range fileList {
		if !file.IsDir() {
			files = append(files, strings.Split(file.Name(), ".")[0])
		}
	}
	return files, err
}
