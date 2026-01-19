---
name: add-to-wishlist
description: Fetch a product page and add it to your wishlist. Extracts product name, type, price, and image from web pages and calls the wishlist CLI tool. Use when adding items to wishlist from URLs or when the user asks to add something to their wishlist.
allowed-tools:
  - WebFetch
  - Bash(/Users/chuhlomin/Projects/chuhlomin/chuhlomin.com/bin/wishlist:*)
  - Read
---

# Add to Wishlist Skill

## Overview

This skill automates adding items to your wishlist YAML files by:
1. Fetching a product web page
2. Extracting product information (name, type, price, image URL)
3. Calling the wishlist CLI tool in programmatic mode to add the item

## Instructions

### Step 1: Understand the Request

Identify the product URL and optionally the wishlist type. The user may say things like:
- "Add https://example.com/product to my wishlist"
- "Add this to my baby wishlist: https://example.com/toy"
- "Put https://shop.com/widget on my watchlist"

Extract:
- **Product URL** (required)
- **Wishlist type** (optional: `default`, `baby`, `watches`, `camera` - defaults to `default`)

### Step 2: Fetch and Extract Product Information

Use the WebFetch tool with a comprehensive extraction prompt:

```
Extract the following product information from this page:
1. Product name: The full product title or name
2. Product type: Categorize this product as one of: book, hardware, toy, wardrobe, font, software, service, subscription, game, accessory, camera, watch, or other relevant category. Infer from the product description, category, and context.
3. Price: The product price with currency symbol (e.g., $50, €200, £75, ₽1000). If multiple prices shown, prefer the main/regular price.
4. Image URL: The primary product image URL or OpenGraph image URL. Must be an absolute URL (starting with http:// or https://).

Format your response clearly with each field labeled.
```

### Step 3: Validate Extracted Data

Ensure all required fields are present:
- **Name**: Must not be empty
- **Type**: Must be a reasonable product category
- **Price**: Must include currency symbol
- **Image URL**: Must be an absolute URL (validate it starts with `http://` or `https://`)

If any required field is missing or invalid:
1. Report which fields are missing or problematic
2. Ask the user if they want to provide the missing information manually

### Step 4: Call the Wishlist CLI Tool

Once all data is validated, call the wishlist CLI in **programmatic mode**:

```bash
/Users/chuhlomin/Projects/chuhlomin/chuhlomin.com/bin/wishlist \
  --url="<product-url>" \
  --name="<extracted-name>" \
  --type="<inferred-type>" \
  --price="<extracted-price>" \
  --image-url="<extracted-image-url>" \
  --list="<wishlist-type>"
```

**Important**: When `--image-url` is provided, the CLI runs in programmatic mode (no interactive prompts).

**Flags**:
- `--url`: Original product page URL (required)
- `--name`: Product name (required in programmatic mode)
- `--type`: Product type/category (required in programmatic mode)
- `--price`: Price with currency (required in programmatic mode)
- `--image-url`: Product image URL (required in programmatic mode)
- `--list`: Wishlist file to update - `default`, `baby`, `watches`, or `camera` (defaults to `default`)

### Step 5: Confirm Success

After the CLI completes:
1. Check the exit code (0 = success)
2. Show the CLI output to the user
3. Confirm which wishlist file was updated (e.g., `content/wishlist.yml`)

If the CLI fails:
1. Show the error message
2. Suggest potential fixes (e.g., rebuild the CLI if binary is missing)

## Examples

### Example 1: Add to Default Wishlist

**User**: "Add https://fsd.it/shop/fonts/pragmatapro/ to my wishlist"

**Actions**:
1. WebFetch: Extract name="Pragmata Pro", type="font", price="€200", image URL
2. Call CLI:
   ```bash
   /Users/chuhlomin/Projects/chuhlomin/chuhlomin.com/bin/wishlist \
     --url="https://fsd.it/shop/fonts/pragmatapro/" \
     --name="Pragmata Pro" \
     --type="font" \
     --price="€200" \
     --image-url="https://fsd.it/path/to/image.jpg" \
     --list="default"
   ```
3. Confirm: "Item added to content/wishlist.yml"

### Example 2: Add to Baby Wishlist

**User**: "Add this toy to my baby wishlist: https://example.com/toy"

**Actions**:
1. WebFetch: Extract product details
2. Call CLI with `--list="baby"`
3. Confirm: "Item added to content/wishlist-baby.yml"

### Example 3: Handle Missing Data

**User**: "Add https://example.com/product"

**Actions**:
1. WebFetch: Name extracted, but price not found on page
2. Report: "I found the product name 'Widget XL' and categorized it as 'hardware', but couldn't find a price on the page. Would you like to provide the price manually?"

## Product Type Categories

When inferring the product type, use these common categories:
- **book** - Books, ebooks, printed materials
- **hardware** - Electronics, gadgets, tools, devices
- **toy** - Toys, games (physical items for children)
- **wardrobe** - Clothing, shoes, accessories, fashion
- **font** - Typefaces, font families
- **software** - Software licenses, apps, programs
- **service** - Subscriptions, memberships, services
- **game** - Video games, board games
- **camera** - Camera equipment, lenses
- **watch** - Watches, timepieces
- **accessory** - General accessories
- **other** - Anything that doesn't fit above categories

## Error Handling

### CLI Binary Missing

If you get "command not found" or similar:
```bash
# Rebuild the wishlist CLI
cd /Users/chuhlomin/Projects/chuhlomin/chuhlomin.com
go build -o bin/wishlist ./cmd/wishlist
```

### WebFetch Fails

If the page can't be fetched:
1. Report the error to the user
2. Suggest they check if the URL is accessible
3. Offer to add the item manually if they provide the details

### Invalid Image URL

If the extracted image URL is relative (e.g., `/images/product.jpg`):
1. Try to resolve it to an absolute URL using the page's base URL
2. If resolution fails, report this and let the CLI handle image fetching (omit `--image-url`)

## Notes

- The wishlist CLI preserves YAML comments when adding items
- Programmatic mode (with `--image-url`) bypasses all interactive prompts
- The CLI also updates the OpenGraph cache at `cache.yml`
- Supported wishlist files:
  - `content/wishlist.yml` (default)
  - `content/wishlist-baby.yml` (baby)
  - `content/wishlist-watches.yml` (watches)
  - `content/wishlist-camera.yml` (camera)

## Tips

- Be generous with type inference - use context clues from the page
- Prefer primary product images over logo/icon images
- If multiple prices exist (sale price, regular price), use the regular price
- Currency symbols are important - don't omit them
- Validate image URLs are absolute before passing to CLI
