// wishlist is a CLI tool for adding items to wishlist YAML files
// with automatic OpenGraph image fetching and DuckDuckGo Image Search fallback.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/chuhlomin/chuhlomin.com/internal/models"
	flags "github.com/jessevdk/go-flags"
)

var (
	cfg Config
)

func main() {
	// Parse flags
	parser := flags.NewParser(&cfg, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatal("Error parsing flags", "error", err)
	}

	// Initialize configuration
	if err := cfg.init(); err != nil {
		log.Fatal("Error initializing config", "error", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatal("Invalid configuration", "error", err)
	}

	// Set log level
	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// Run the main workflow
	if err := run(); err != nil {
		log.Fatal("Error", "error", err)
	}
}

func run() error {
	log.Info("Starting wishlist CLI")
	log.Debugf("Configuration: %+v", cfg)

	// Create OpenGraph client with 5s timeout
	ogClient, err := newOpenGraphClient(5*time.Second, cfg.CacheFile)
	if err != nil {
		return fmt.Errorf("failed to create OpenGraph client: %w", err)
	}
	defer func() {
		if err := ogClient.Save(cfg.CacheFile); err != nil {
			log.Error("Failed to save cache", "error", err)
		}
	}()

	// Create wishlist manager
	wlManager := NewWishlistManager(cfg.WishlistFile)

	// Determine image URL
	var imageURL string

	if cfg.ImageURL != "" {
		// Programmatic mode: use provided image URL
		imageURL = cfg.ImageURL
		log.Infof("Using provided image URL: %s", imageURL)
	} else {
		// Try to fetch OpenGraph metadata
		log.Infof("[1/5] Fetching OpenGraph metadata from %s...", cfg.ItemURL)
		og, err := ogClient.Get(cfg.ItemURL)

		var needsImageSearch bool

		if err == nil && og.Image != "" {
			// Resolve relative URLs to absolute
			resolvedURL, resolveErr := resolveURL(cfg.ItemURL, og.Image)
			if resolveErr != nil {
				log.Warnf("Failed to resolve image URL: %v", resolveErr)
				needsImageSearch = true
			} else {
				log.Infof("✓ Found OpenGraph image: %s", resolvedURL)

				// Display the OG image for confirmation
				display := NewImageDisplay()
				fmt.Println("\n=== OpenGraph Image ===")
				if err := display.DisplayImage(resolvedURL, 1); err != nil {
					log.Warnf("Failed to display image: %v", err)
				}

				// Ask user if they want to use this image
				confirmed, confirmErr := ConfirmImage(resolvedURL)
				if confirmErr != nil {
					log.Warnf("Failed to get confirmation: %v", confirmErr)
					needsImageSearch = true
				} else if confirmed {
					// User confirmed, use this image
					imageURL = resolvedURL
					log.Info("✓ Using OpenGraph image")
				} else {
					// User declined, search for alternatives
					log.Info("OpenGraph image declined, searching for alternatives...")
					needsImageSearch = true
				}
			}
		} else {
			// No OpenGraph image found
			if err != nil {
				log.Warnf("✗ Failed to fetch OpenGraph data: %v", err)
			} else {
				log.Warn("✗ No OpenGraph image found")
			}
			needsImageSearch = true
		}

		// Search DuckDuckGo if needed
		if needsImageSearch && imageURL == "" {
			log.Info("[2/5] Searching for images with DuckDuckGo...")
			ddgClient := NewDuckDuckGoImageSearchClient()

			// Use item name as search query, or URL host if name not provided
			searchQuery := cfg.ItemName
			if searchQuery == "" {
				// Extract domain from URL for search
				searchQuery = cfg.ItemURL
			}

			results, err := ddgClient.Search(searchQuery, 4)
			if err != nil {
				log.Warnf("✗ DuckDuckGo search failed: %v", err)
				log.Info("Please provide image URL manually (--image-url)")
			} else if len(results) == 0 {
				log.Warn("✗ No images found")
				log.Info("Please provide image URL manually (--image-url)")
			} else {
				log.Infof("✓ Found %d images", len(results))

				// Display images with Kitty graphics protocol
				display := NewImageDisplay()
				if err := display.DisplayImages(results); err != nil {
					log.Warnf("Failed to display images: %v", err)
				}

				// Let user select an image
				selectedURL, err := SelectImage(results)
				if err != nil {
					log.Errorf("Error selecting image: %v", err)
					return fmt.Errorf("image selection failed: %w", err)
				}
				imageURL = selectedURL
				if imageURL != "" {
					log.Infof("Selected image: %s", imageURL)
				} else {
					log.Info("No image selected")
				}
			}
		}
	}

	// Update cache if we have an image
	if imageURL != "" {
		ogClient.Set(cfg.ItemURL, imageURL)
	}

	// Collect item details
	var itemName, itemType, itemPrice string

	// In programmatic mode, all fields are required
	// In interactive mode, prompt for missing fields
	if cfg.IsProgrammaticMode() {
		itemName = cfg.ItemName
		itemType = cfg.ItemType
		itemPrice = cfg.ItemPrice
	} else {
		log.Info("[4/5] Collecting item details...")

		// Prompt for name
		var promptErr error
		itemName, promptErr = PromptForName(cfg.ItemName)
		if promptErr != nil {
			return fmt.Errorf("failed to get item name: %w", promptErr)
		}

		// Prompt for type
		if cfg.ItemType != "" {
			itemType = cfg.ItemType
		} else {
			itemType, promptErr = PromptForType()
			if promptErr != nil {
				return fmt.Errorf("failed to get item type: %w", promptErr)
			}
		}

		// Prompt for price
		if cfg.ItemPrice != "" {
			itemPrice = cfg.ItemPrice
		} else {
			itemPrice, promptErr = PromptForPrice()
			if promptErr != nil {
				return fmt.Errorf("failed to get item price: %w", promptErr)
			}
		}
	}

	// Create wishlist item
	item := models.WishlistItem{
		Name:  itemName,
		Type:  itemType,
		URL:   cfg.ItemURL,
		Price: itemPrice,
	}

	// Confirm before adding (only in interactive mode)
	if !cfg.IsProgrammaticMode() {
		confirmed, confirmErr := ConfirmItem(item, imageURL)
		if confirmErr != nil {
			return fmt.Errorf("failed to get confirmation: %w", confirmErr)
		}
		if !confirmed {
			log.Info("Operation cancelled by user")
			return nil
		}
	} else {
		// Show summary in programmatic mode
		fmt.Println("\n" + "=== Summary ===")
		fmt.Printf("Name:  %s\n", item.Name)
		fmt.Printf("Type:  %s\n", item.Type)
		fmt.Printf("URL:   %s\n", item.URL)
		fmt.Printf("Price: %s\n", item.Price)
		fmt.Printf("Image: %s\n", imageURL)
	}

	// Add to wishlist
	log.Infof("[5/5] Adding item to %s...", cfg.WishlistFile)
	if err := wlManager.AppendItem(item); err != nil {
		return fmt.Errorf("failed to add item to wishlist: %w", err)
	}
	log.Info("✓ Item added to wishlist")

	fmt.Println("\n✓ Done!")
	fmt.Printf("  - Added to %s\n", cfg.WishlistFile)
	fmt.Printf("  - Updated %s\n", cfg.CacheFile)

	return nil
}
