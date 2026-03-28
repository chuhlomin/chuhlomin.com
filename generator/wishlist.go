package main

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"

	"github.com/chuhlomin/chuhlomin.com/internal/models"
)

// Exchange rates to USD (approximate)
var exchangeRates = map[string]float64{
	"$": 1.0,    // USD
	"€": 1.08,   // EUR
	"£": 1.26,   // GBP
	"₽": 0.011,  // RUB
}

// parsePrice parses a price string and returns the price in USD cents.
// Returns -1 for unknown or unparseable prices.
func parsePrice(price string) int {
	price = strings.TrimSpace(price)
	if price == "" || price == "(unknown)" {
		return -1
	}

	// Try to find currency symbol and extract number
	var currency string
	var rate float64

	for symbol, r := range exchangeRates {
		if strings.Contains(price, symbol) {
			currency = symbol
			rate = r
			break
		}
	}

	if currency == "" {
		return -1
	}

	// Extract numeric value using regex (handles formats like $1,300 or 490 ₽)
	re := regexp.MustCompile(`[\d,]+(?:\.\d+)?`)
	match := re.FindString(price)
	if match == "" {
		return -1
	}

	// Remove commas and parse
	match = strings.ReplaceAll(match, ",", "")
	value, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return -1
	}

	// Convert to USD cents
	usdCents := int(value * rate * 100)
	return usdCents
}

type WishlistFile struct {
	Language string
}

type PageData struct {
	File  WishlistFile
	Title string
	Items interface{}
}

func (g *Generator) processWishlistItems(fileContent []byte, filename string) (PageData, error) {
	// unmarshal the file into array of WishlistItem
	var items []models.WishlistItem
	err := yaml.Unmarshal(fileContent, &items)
	if err != nil {
		return PageData{}, fmt.Errorf("Error unmarshaling file: %v", err)
	}

	// reverse the items, so that the newest ones are first
	slices.Reverse(items)

	// enrich the items with OpenGraph data and parse prices
	for i := range items {
		og, err := g.og.Get(items[i].URL)
		if err != nil {
			log.Errorf("Error getting OpenGraph data for %s: %v", items[i].URL, err)
			continue
		}

		items[i].Image = og.Image
		items[i].NormalizedPrice = parsePrice(items[i].Price)
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
