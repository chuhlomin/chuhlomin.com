package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/BourgeoisBear/rasterm"
)

// ImageDisplay handles displaying images in the terminal
type ImageDisplay struct {
	supportsKitty bool
	supportsIterm bool
	supportsSixel bool
}

// NewImageDisplay creates a new image display handler
func NewImageDisplay() *ImageDisplay {
	// Detect terminal capabilities
	// rasterm automatically detects and uses the best protocol available
	supportsKitty := rasterm.IsKittyCapable()
	supportsIterm := rasterm.IsItermCapable()
	supportsSixel, _ := rasterm.IsSixelCapable()

	return &ImageDisplay{
		supportsKitty: supportsKitty,
		supportsIterm: supportsIterm,
		supportsSixel: supportsSixel,
	}
}

// SupportsGraphics returns true if terminal supports any graphics protocol
func (d *ImageDisplay) SupportsGraphics() bool {
	return d.supportsKitty || d.supportsIterm || d.supportsSixel
}

// DisplayImage displays an image from a URL in the terminal
// Returns error if image cannot be fetched or displayed
func (d *ImageDisplay) DisplayImage(imageURL string, index int) error {
	// Fetch image
	img, err := fetchImage(imageURL)
	if err != nil {
		// Fallback to URL on error
		fmt.Printf("  %d. %s (failed to load: %v)\n", index, imageURL, err)
		return nil
	}

	fmt.Printf("\n  Image %d:\n", index)

	// Try to display using available terminal protocol
	var displayErr error
	if d.supportsKitty {
		// Use Kitty protocol (preferred)
		displayErr = rasterm.KittyWriteImage(os.Stdout, img, rasterm.KittyImgOpts{})
	} else if d.supportsIterm {
		// Use iTerm2 protocol
		displayErr = rasterm.ItermWriteImage(os.Stdout, img)
	} else if d.supportsSixel {
		// Sixel requires paletted image - for simplicity, skip for now
		// In production, we'd convert the image to paletted format
		fmt.Printf("  %s (Sixel not fully implemented)\n", imageURL)
		return nil
	} else {
		// No graphics support - just show URL
		fmt.Printf("  %s\n", imageURL)
		return nil
	}

	if displayErr != nil {
		// Fallback to URL on error
		fmt.Printf("  %s (failed to display: %v)\n", imageURL, displayErr)
		return nil
	}

	fmt.Printf("  %s\n\n", imageURL)
	return nil
}

// DisplayImages displays multiple images with labels
func (d *ImageDisplay) DisplayImages(images []ImageSearchResult) error {
	if len(images) == 0 {
		return fmt.Errorf("no images to display")
	}

	fmt.Println("\n=== Available Images ===")

	for i, img := range images {
		if err := d.DisplayImage(img.ImageURL, i+1); err != nil {
			return err
		}
	}

	return nil
}

// fetchImage downloads an image from a URL and decodes it
func fetchImage(imageURL string) (image.Image, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Decode image
	img, _, err := image.Decode(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, nil
}
