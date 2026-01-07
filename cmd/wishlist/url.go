package main

import (
	"net/url"
	"strings"
)

// resolveURL resolves a potentially relative URL against a base URL
// If imageURL is already absolute, returns it unchanged
// If imageURL is relative, combines it with the base URL
func resolveURL(baseURL, imageURL string) (string, error) {
	// If image URL is already absolute, return as-is
	if strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://") {
		return imageURL, nil
	}

	// Parse base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// Parse relative URL
	ref, err := url.Parse(imageURL)
	if err != nil {
		return "", err
	}

	// Resolve relative URL against base
	resolved := base.ResolveReference(ref)
	return resolved.String(), nil
}
