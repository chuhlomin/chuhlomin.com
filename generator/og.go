package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"gopkg.in/yaml.v3"
)

type OpenGraph struct {
	// Title       string `yaml:"title,omitempty"`
	// Description string `yaml:"description,omitempty"`
	Image string `yaml:"image,omitempty"`
	// ImageWidth  int    `yaml:"image_width,omitempty"`
	// ImageHeight int    `yaml:"image_height,omitempty"`
}

type openGraphClient struct {
	client *http.Client
	cache  map[string]OpenGraph
}

func newOpenGraphClient(clientTimeout time.Duration, cacheFile string) (*openGraphClient, error) {
	cache, err := readCache(cacheFile)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: clientTimeout,
	}

	return &openGraphClient{client: client, cache: cache}, nil
}

func (c *openGraphClient) Get(url string) (OpenGraph, error) {
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

func (c *openGraphClient) Save(cacheFile string) error {
	f, err := os.Create(cacheFile)
	if err != nil {
		return err
	}

	defer f.Close()

	e := yaml.NewEncoder(f)
	e.SetIndent(2)
	return e.Encode(c.cache)
}

func (c *openGraphClient) fetch(url string) (OpenGraph, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return OpenGraph{}, err
	}

	defer resp.Body.Close()

	return parseOpenGraph(resp.Body)
}

func readCache(cacheFile string) (cache map[string]OpenGraph, err error) {
	f, err := os.Open(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			cache = make(map[string]OpenGraph)
			return cache, nil
		}

		return nil, err
	}

	defer f.Close()
	err = yaml.NewDecoder(f).Decode(&cache)
	if err != nil {
		return nil, err
	}
	return cache, nil
}

func parseOpenGraph(r io.Reader) (og OpenGraph, err error) {
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
