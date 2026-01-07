package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/chuhlomin/chuhlomin.com/internal/cache"
)

// openGraphClient is a client for fetching OpenGraph metadata with caching
type openGraphClient struct {
	client *http.Client
	cache  map[string]cache.OpenGraph
}

// newOpenGraphClient creates a new OpenGraph client with the given timeout and cache file
func newOpenGraphClient(clientTimeout time.Duration, cacheFile string) (*openGraphClient, error) {
	ogCache, err := cache.ReadCache(cacheFile)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: clientTimeout,
	}

	return &openGraphClient{client: client, cache: ogCache}, nil
}

// Get retrieves OpenGraph metadata for the given URL (from cache or by fetching)
func (c *openGraphClient) Get(url string) (cache.OpenGraph, error) {
	og, ok := c.cache[url]
	if ok {
		return og, nil
	}

	og, err := c.fetch(url)
	if err != nil {
		return og, err
	}

	c.cache[url] = og

	return og, nil
}

// Save writes the cache to disk
func (c *openGraphClient) Save(cacheFile string) error {
	return cache.SaveCache(cacheFile, c.cache)
}

// Set adds or updates an OpenGraph entry in the cache
func (c *openGraphClient) Set(url, imageURL string) {
	c.cache[url] = cache.OpenGraph{Image: imageURL}
}

// fetch retrieves OpenGraph metadata from the given URL
func (c *openGraphClient) fetch(url string) (cache.OpenGraph, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return cache.OpenGraph{}, err
	}

	defer resp.Body.Close()

	return parseOpenGraph(resp.Body)
}

// parseOpenGraph extracts OpenGraph image from HTML content
func parseOpenGraph(r io.Reader) (og cache.OpenGraph, err error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return og, err
	}

	// extract "og:image" tag
	re := regexp.MustCompile(`<meta property="og:image" content="([^"]+)"`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) != 2 {
		return og, fmt.Errorf("could not find og:image tag")
	}

	og.Image = matches[1]

	return og, nil
}
