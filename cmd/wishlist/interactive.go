package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chuhlomin/chuhlomin.com/internal/models"
)

// PromptForName prompts the user to enter an item name
func PromptForName(defaultName string) (string, error) {
	return promptWithDefault("Enter item name", defaultName)
}

// PromptForType prompts the user to enter an item type
func PromptForType() (string, error) {
	fmt.Println("\nCommon types: book, hardware, wardrobe, furniture, thing, toy, kitchen, health, travel, font")
	return prompt("Enter item type")
}

// PromptForPrice prompts the user to enter an item price
func PromptForPrice() (string, error) {
	fmt.Println("\nExamples: $20, €50, ₽490, £30")
	return prompt("Enter item price")
}

// SelectImage prompts the user to select an image from a list of results
// Returns the selected image URL
func SelectImage(images []ImageSearchResult) (string, error) {
	if len(images) == 0 {
		return "", fmt.Errorf("no images to select from")
	}

	fmt.Println("\nSelect an image:")
	fmt.Printf("  Enter 1-%d to select an image\n", len(images))
	fmt.Println("  Or enter a custom image URL")
	fmt.Println("  Or press Enter to skip (no image)")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nYour choice: ")

	if !scanner.Scan() {
		return "", fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(scanner.Text())

	// Empty input means skip
	if input == "" {
		return "", nil
	}

	// Try to parse as a number (1-based index)
	if index, err := strconv.Atoi(input); err == nil {
		if index < 1 || index > len(images) {
			return "", fmt.Errorf("invalid selection: must be between 1 and %d", len(images))
		}
		return images[index-1].ImageURL, nil
	}

	// Otherwise, treat as custom URL
	// Basic URL validation
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		return "", fmt.Errorf("invalid URL: must start with http:// or https://")
	}

	return input, nil
}

// ConfirmImage shows an image and asks if user wants to use it
func ConfirmImage(imageURL string) (bool, error) {
	return promptYesNo("Use this image?")
}

// ConfirmItem shows the item summary and asks for confirmation
func ConfirmItem(item models.WishlistItem, imageURL string) (bool, error) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("=== Summary ===")
	fmt.Printf("Name:  %s\n", item.Name)
	fmt.Printf("Type:  %s\n", item.Type)
	fmt.Printf("URL:   %s\n", item.URL)
	fmt.Printf("Price: %s\n", item.Price)
	fmt.Printf("Image: %s\n", imageURL)
	fmt.Println(strings.Repeat("=", 50))

	return promptYesNo("Add this item to the wishlist?")
}

// prompt asks the user for input with a prompt message
func prompt(message string) (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s: ", message)

	if !scanner.Scan() {
		return "", fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return "", fmt.Errorf("input cannot be empty")
	}

	return input, nil
}

// promptWithDefault asks the user for input with a default value
func promptWithDefault(message, defaultValue string) (string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", message, defaultValue)
	} else {
		fmt.Printf("%s: ", message)
	}

	if !scanner.Scan() {
		return "", fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(scanner.Text())

	// Use default if input is empty and default is provided
	if input == "" && defaultValue != "" {
		return defaultValue, nil
	}

	if input == "" {
		return "", fmt.Errorf("input cannot be empty")
	}

	return input, nil
}

// promptYesNo asks a yes/no question and returns true for yes
func promptYesNo(message string) (bool, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s (y/n): ", message)

	if !scanner.Scan() {
		return false, fmt.Errorf("failed to read input")
	}

	input := strings.ToLower(strings.TrimSpace(scanner.Text()))

	switch input {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid input: please enter 'y' or 'n'")
	}
}
