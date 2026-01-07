package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

// ImageSearchResult represents a single image result from image search
type ImageSearchResult struct {
	Title        string
	ThumbnailURL string
	ImageURL     string
	Width        int
	Height       int
}

// DuckDuckGoImageSearchClient is a client for DuckDuckGo Image Search
type DuckDuckGoImageSearchClient struct {
	HTTPClient *http.Client
}

// NewDuckDuckGoImageSearchClient creates a new DuckDuckGo Image Search client
func NewDuckDuckGoImageSearchClient() *DuckDuckGoImageSearchClient {
	return &DuckDuckGoImageSearchClient{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ddgImageResult represents the internal DuckDuckGo API response structure
type ddgImageResult struct {
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	Image     string `json:"image"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// ddgImageResponse represents the DuckDuckGo image search API response
type ddgImageResponse struct {
	Results []ddgImageResult `json:"results"`
}

// Search searches for images using DuckDuckGo
// Returns up to 'count' image results for the given query
func (c *DuckDuckGoImageSearchClient) Search(query string, count int) ([]ImageSearchResult, error) {
	// Step 1: Get the vqd token
	vqd, err := c.getVQDToken(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get VQD token: %w", err)
	}

	// Step 2: Search for images using the vqd token
	apiURL := "https://duckduckgo.com/i.js"
	params := url.Values{}
	params.Set("q", query)
	params.Set("o", "json")
	params.Set("vqd", vqd)

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	// Create request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://duckduckgo.com/")

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s (body: %s)", resp.StatusCode, resp.Status, string(body))
	}

	// Parse response
	var searchResp ddgImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our result type
	results := make([]ImageSearchResult, 0, len(searchResp.Results))
	for i, r := range searchResp.Results {
		if i >= count {
			break
		}
		results = append(results, ImageSearchResult{
			Title:        r.Title,
			ThumbnailURL: r.Thumbnail,
			ImageURL:     r.Image,
			Width:        r.Width,
			Height:       r.Height,
		})
	}

	return results, nil
}

// getVQDToken fetches the vqd token required for DuckDuckGo image search
func (c *DuckDuckGoImageSearchClient) getVQDToken(query string) (string, error) {
	// Make a request to the main search page to get the vqd token
	searchURL := fmt.Sprintf("https://duckduckgo.com/?q=%s&iax=images&ia=images", url.QueryEscape(query))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Extract vqd token from the response using regex
	// The vqd token appears in the HTML like: vqd="3-..."
	re := regexp.MustCompile(`vqd=['"]([^'"]+)['"]`)
	matches := re.FindSubmatch(body)
	if len(matches) < 2 {
		return "", fmt.Errorf("vqd token not found in response")
	}

	return string(matches[1]), nil
}
