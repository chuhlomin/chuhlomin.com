// Package models contains shared data models used across the project
package models

// WishlistItem represents an item in a wishlist YAML file
type WishlistItem struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	URL   string `yaml:"url"`
	Price string `yaml:"price"`
	Image string `yaml:"image,omitempty"` // OpenGraph image URL (populated by generator/CLI)
}
