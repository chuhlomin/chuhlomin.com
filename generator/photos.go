package main

import (
	"fmt"
	"slices"

	"gopkg.in/yaml.v3"
)

// Photos is a struct describing photos.yml file
type Photos struct {
	Root   string
	Photos []Photo
}

// Photo is a struct for items in photos.yml file
type Photo struct {
	Path string
}

func (g *Generator) processPhotos(fileContent []byte) (interface{}, error) {
	var photos Photos
	err := yaml.Unmarshal(fileContent, &photos)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling file: %v", err)
	}

	slices.Reverse(photos.Photos)

	return photos, nil
}
