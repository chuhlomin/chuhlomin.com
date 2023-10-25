package main

import (
	"fmt"
	"slices"

	"gopkg.in/yaml.v3"
)

// Photo is a struct for items in photos.yml file
type Photo struct {
	Path string
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
