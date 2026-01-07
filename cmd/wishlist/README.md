# Wishlist CLI

A command-line tool for adding items to wishlist YAML files with automatic OpenGraph image fetching, DuckDuckGo Image Search fallback, and interactive image selection using the Kitty graphics protocol.

## Features

- ✓ Automatic OpenGraph metadata extraction from product URLs
- ✓ Relative URL resolution (prepends base URL to relative image paths)
- ✓ Interactive image confirmation with in-terminal preview (Kitty graphics)
- ✓ DuckDuckGo Image Search fallback when OpenGraph declined or not found (no API key required)
- ✓ Interactive prompts for missing item details
- ✓ Support for multiple wishlist files (default, baby, watches, camera)
- ✓ Programmatic mode for automation (skip all prompts)
- ✓ Cache system for OpenGraph images (shared with generator)

## Installation

```bash
# From repository root
cd cmd/wishlist
go build -o ../../bin/wishlist

# Or use from anywhere in the repository
go build -C cmd/wishlist -o bin/wishlist
```

## Usage

### Interactive Mode (Default)

The simplest way to add an item:

```bash
wishlist --url "https://example.com/product"
```

This will:
1. Fetch OpenGraph image from the product URL
2. If found, display the image in terminal and ask for confirmation
3. If you decline the OG image (or none found), search DuckDuckGo for alternatives
4. Prompt for missing details (name, type, price)
5. Show final confirmation before adding to wishlist

### Programmatic Mode

Skip all prompts by providing all required fields:

```bash
wishlist \
  --url "https://example.com/product" \
  --name "Product Name" \
  --type "hardware" \
  --price "$99" \
  --image-url "https://example.com/image.jpg"
```

### Specify Wishlist File

Use the `--list` flag to choose which wishlist to update:

```bash
# Add to baby wishlist
wishlist --url "https://example.com/toy" --list baby

# Add to watches wishlist
wishlist --url "https://example.com/watch" --list watches

# Add to camera wishlist
wishlist --url "https://example.com/lens" --list camera

# Add to default wishlist (default)
wishlist --url "https://example.com/item"
```

## Command-Line Flags

| Flag | Description | Required | Default |
|------|-------------|----------|---------|
| `--url` | Item URL | Yes | - |
| `--name` | Item name | No* | Prompts if missing |
| `--type` | Item type (e.g., book, hardware, toy) | No* | Prompts if missing |
| `--price` | Item price with currency (e.g., $20, €50) | No* | Prompts if missing |
| `--list` | Wishlist type: default, baby, watches, camera | No | default |
| `--image-url` | Direct image URL (skips OG/DuckDuckGo search) | No | - |
| `--debug` | Enable debug logging | No | false |

\* Required when using `--image-url` (programmatic mode)

## Environment Variables

No environment variables are required. The tool uses DuckDuckGo Image Search which does not require an API key.

## Configuration

The tool uses the following paths (relative to working directory):

- **Content directory**: `content/`
- **Wishlist files**:
  - `content/wishlist.yml` (default)
  - `content/wishlist-baby.yml` (baby)
  - `content/wishlist-watches.yml` (watches)
  - `content/wishlist-camera.yml` (camera)
- **Cache file**: `cache.yml`

## Item Types

Common item types used in wishlists:

- `book`, `font` - Reading materials and typography
- `hardware`, `software` - Tech products
- `wardrobe`, `furniture` - Home and personal items
- `thing`, `toy` - General items and toys
- `kitchen`, `health`, `travel` - Lifestyle categories

## Examples

### Example 1: Full Interactive Flow (Accept OG Image)

```bash
$ wishlist --url "https://pragmatapro.com/"

[1/5] Fetching OpenGraph metadata from https://pragmatapro.com/...
✓ Found OpenGraph image: https://pragmatapro.com/og-image.jpg

=== OpenGraph Image ===
  Image 1:
  [Kitty displays image preview]
  https://pragmatapro.com/og-image.jpg

Use this image? (y/n): y
✓ Using OpenGraph image

[4/5] Collecting item details...
Enter item name: Pragmata Pro

Common types: book, hardware, wardrobe, furniture, thing, toy, kitchen, health, travel, font
Enter item type: font

Examples: $20, €50, ₽490, £30
Enter item price: €200

==================================================
=== Summary ===
Name:  Pragmata Pro
Type:  font
URL:   https://pragmatapro.com/
Price: €200
Image: https://pragmatapro.com/og-image.jpg
==================================================
Add this item to the wishlist? (y/n): y

[5/5] Adding item to content/wishlist.yml...
✓ Item added to wishlist

✓ Done!
  - Added to content/wishlist.yml
  - Updated cache.yml
```

### Example 2: Decline OG Image → DuckDuckGo Search

```bash
$ wishlist --url "https://example.com/product"

[1/5] Fetching OpenGraph metadata...
✓ Found OpenGraph image: https://example.com/og-image.jpg

=== OpenGraph Image ===
  Image 1:
  [Kitty displays image preview]
  https://example.com/og-image.jpg

Use this image? (y/n): n
OpenGraph image declined, searching for alternatives...

[2/5] Searching for images with DuckDuckGo...
✓ Found 4 images

=== Available Images ===
(Kitty terminal not detected - showing URLs only)
Tip: Use Kitty terminal to see image previews

  1. https://example.com/img1.jpg
  2. https://example.com/img2.jpg
  3. https://example.com/img3.jpg
  4. https://example.com/img4.jpg

Select an image:
  Enter 1-4 to select an image
  Or enter a custom image URL
  Or press Enter to skip (no image)

Your choice: 2
Selected image: https://example.com/img2.jpg

[4/5] Collecting item details...
...
```

### Example 3: Programmatic Mode (Automation)

```bash
wishlist \
  --url "https://example.com/product" \
  --name "Example Product" \
  --type "hardware" \
  --price "$49.99" \
  --image-url "https://example.com/product.jpg" \
  --list baby

=== Summary ===
Name:  Example Product
Type:  hardware
URL:   https://example.com/product
Price: $49.99
Image: https://example.com/product.jpg

[5/5] Adding item to content/wishlist-baby.yml...
✓ Item added to wishlist

✓ Done!
  - Added to content/wishlist-baby.yml
  - Updated cache.yml
```

## Kitty Terminal Support

For the best experience with image previews, use [Kitty terminal](https://sw.kovidgoyal.net/kitty/):

```bash
# macOS
brew install kitty

# Then run wishlist from Kitty terminal
```

When running in a non-Kitty terminal, the tool gracefully falls back to showing image URLs instead of previews.

## Integration with Other Apps

The programmatic mode (`--image-url`) allows other applications to call this tool:

```bash
# Another app can pass an image URL directly
other-app --get-image | xargs -I {} wishlist \
  --url "$URL" \
  --name "$NAME" \
  --type "$TYPE" \
  --price "$PRICE" \
  --image-url "{}"
```

## Relative URL Handling

When OpenGraph returns a relative image URL (e.g., `/images/product.jpg`), the tool automatically resolves it to an absolute URL by prepending the base URL:

```
Item URL:      https://example.com/products/item
OG Image:      /images/product.jpg
Resolved to:   https://example.com/images/product.jpg
```

This is saved in `cache.yml` as the absolute URL for future use.

## Files Modified

When you run the tool, it updates:

1. **Wishlist YAML file** (`content/wishlist*.yml`):
   - Appends new item to the end (generator reverses order when building site)
   - Preserves existing items
   - Uses 2-space YAML indentation

2. **Cache file** (`cache.yml`):
   - Stores URL → image mappings (always as absolute URLs)
   - Shared with the generator to avoid re-fetching OpenGraph data
   - Persists across runs

## Troubleshooting

### "Failed to fetch OpenGraph data: timeout"

The website may be slow or blocking automated requests. Options:

1. Use `--image-url` to provide the image directly
2. Wait and try again later
3. Check if the website requires authentication

### "No images found" from DuckDuckGo

The search query might be too specific. Try:

1. Providing a `--name` flag with a more general product name
2. Using `--image-url` to provide the image directly

### Images not displaying in terminal

Make sure you're using Kitty terminal. Check with:

```bash
echo $TERM  # Should show "xterm-kitty"
```

## Development

### Building

```bash
go build -o ../../bin/wishlist
```

### Testing

```bash
# Test with debug logging
wishlist --url "https://example.com" --debug

# Test different wishlist files
wishlist --url "https://example.com" --list baby
wishlist --url "https://example.com" --list watches
wishlist --url "https://example.com" --list camera
```

### Dependencies

- `github.com/charmbracelet/log` - Structured logging
- `github.com/jessevdk/go-flags` - CLI flag parsing
- `github.com/dolmen-go/kittyimg` - Kitty graphics protocol
- `gopkg.in/yaml.v3` - YAML parsing and marshaling

## License

Part of the chuhlomin.com repository.
