package main

import (
	"fmt"
	"slices"

	"gopkg.in/yaml.v3"
)

// Photo is a struct for items in photos.yml file
type Photo struct {
	Path                string
	Title               string `yaml:"title,omitempty"`
	TitleRu             string `yaml:"title_ru,omitempty"`
	ThumbPath           string `yaml:"thumb,omitempty"`
	Blurhash            string `yaml:"blurhash,omitempty"`
	BlurhashImageBase64 string `yaml:"blurhash_image_base64,omitempty"`
	Width               int    `yaml:"width,omitempty"`
	Height              int    `yaml:"height,omitempty"`
	ThumbXOffset        int    `yaml:"thumb_x,omitempty"`
	ThumbYOffset        int    `yaml:"thumb_y,omitempty"`
	ThumbWidth          int    `yaml:"thumb_width,omitempty"`
	ThumbHeight         int    `yaml:"thumb_height,omitempty"`
	ThumbTotalWidth     int    `yaml:"thumb_total_width,omitempty"`
	ThumbTotalHeight    int    `yaml:"thumb_total_height,omitempty"`
}

func (g *Generator) processPhotos(fileContent []byte) (interface{}, error) {
	var photos []Photo
	err := yaml.Unmarshal(fileContent, &photos)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling file: %v", err)
	}

	slices.Reverse(photos)

	return photos, nil
}
