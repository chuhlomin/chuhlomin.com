// Package cache provides OpenGraph metadata caching functionality
package cache

import (
	"os"

	"gopkg.in/yaml.v3"
)

// OpenGraph holds OpenGraph metadata extracted from a webpage
type OpenGraph struct {
	Image string `yaml:"image,omitempty"`
}

// ReadCache loads the OpenGraph cache from disk
// Returns an empty cache if the file doesn't exist
func ReadCache(cacheFile string) (map[string]OpenGraph, error) {
	f, err := os.Open(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]OpenGraph), nil
		}
		return nil, err
	}
	defer f.Close()

	var cache map[string]OpenGraph
	err = yaml.NewDecoder(f).Decode(&cache)
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// SaveCache writes the OpenGraph cache to disk with 2-space YAML indentation
func SaveCache(cacheFile string, cache map[string]OpenGraph) error {
	f, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer f.Close()

	e := yaml.NewEncoder(f)
	e.SetIndent(2)
	return e.Encode(cache)
}
