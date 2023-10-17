package main

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

// WishlistItem is a struct for items in wishlist.yml file
type WishlistItem struct {
	Name  string
	Type  string
	URL   string
	Price string

	Image string // OpenGraph image URL
}

func (g *Generator) processWishlistItems(fileContent []byte) (interface{}, error) {
	// unmarshal the file into array of WishlistItem
	var items []WishlistItem
	err := yaml.Unmarshal(fileContent, &items)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling file: %v", err)
	}

	// reverse the items, so that the newest ones are first
	slices.Reverse(items)

	// enrich the items with OpenGraph data
	for i := range items {
		og, err := g.og.Get(items[i].URL)
		if err != nil {
			log.Errorf("Error getting OpenGraph data for %s: %v", items[i].URL, err)
			continue
		}

		items[i].Image = og.Image
	}

	return items, nil
}
