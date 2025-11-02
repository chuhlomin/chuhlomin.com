package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type WishlistFile struct {
	Language string
}

type PageData struct {
	File  WishlistFile
	Title string
	Items interface{}
}

// WishlistItem is a struct for items in wishlist.yml file
type WishlistItem struct {
	Name  string
	Type  string
	URL   string
	Price string

	Image string // OpenGraph image URL
}

func (g *Generator) processWishlistItems(fileContent []byte, filename string) (PageData, error) {
	// unmarshal the file into array of WishlistItem
	var items []WishlistItem
	err := yaml.Unmarshal(fileContent, &items)
	if err != nil {
		return PageData{}, fmt.Errorf("Error unmarshaling file: %v", err)
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

	return PageData{
		File: WishlistFile{
			Language: "en",
		},
		Title: generateTitleFromFilename(filename),
		Items: items,
	}, nil
}

func generateTitleFromFilename(filename string) string {
	caser := cases.Title(language.English, cases.NoLower)

	return caser.String(
		strings.ReplaceAll(
			strings.TrimSuffix(filename, ".yml"),
			"-",
			" ",
		),
	)
}
