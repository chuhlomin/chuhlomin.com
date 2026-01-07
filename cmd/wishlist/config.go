package main

import (
	"fmt"
	"path/filepath"
)

// Config holds all configuration options for the wishlist CLI
type Config struct {
	// CLI flags
	ItemURL      string `long:"url" description:"Item URL" required:"true"`
	ItemName     string `long:"name" description:"Item name (optional, will prompt if missing)"`
	ItemType     string `long:"type" description:"Item type (e.g., book, hardware, toy)"`
	ItemPrice    string `long:"price" description:"Item price with currency (e.g., $20, €50, ₽490)"`
	WishlistType string `long:"list" description:"Wishlist type" choice:"default" choice:"baby" choice:"watches" choice:"camera" default:"default"`
	ImageURL     string `long:"image-url" description:"Direct image URL (skips OpenGraph fetch and DuckDuckGo search)"`

	// Derived paths (computed in init())
	WishlistFile string
	CacheFile    string
	ContentDir   string

	// Debug mode
	Debug bool `long:"debug" description:"Enable debug logging"`
}

// init computes derived paths based on configuration
func (c *Config) init() error {
	// Set content directory (relative to working directory)
	if c.ContentDir == "" {
		c.ContentDir = "content"
	}

	// Compute wishlist file path based on type
	c.WishlistFile = c.getWishlistFilePath()

	// Set cache file path
	if c.CacheFile == "" {
		c.CacheFile = "cache.yml"
	}

	return nil
}

// getWishlistFilePath returns the path to the wishlist YAML file based on type
func (c *Config) getWishlistFilePath() string {
	var filename string

	switch c.WishlistType {
	case "baby":
		filename = "wishlist-baby.yml"
	case "watches":
		filename = "wishlist-watches.yml"
	case "camera":
		filename = "wishlist-camera.yml"
	default:
		filename = "wishlist.yml"
	}

	return filepath.Join(c.ContentDir, filename)
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// If --image-url is provided, name, type, and price must also be provided
	if c.ImageURL != "" {
		if c.ItemName == "" {
			return fmt.Errorf("--name is required when using --image-url")
		}
		if c.ItemType == "" {
			return fmt.Errorf("--type is required when using --image-url")
		}
		if c.ItemPrice == "" {
			return fmt.Errorf("--price is required when using --image-url")
		}
	}

	return nil
}

// IsProgrammaticMode returns true if the tool is being run in programmatic mode
// (i.e., --image-url is provided, skipping interactive prompts)
func (c *Config) IsProgrammaticMode() bool {
	return c.ImageURL != ""
}
