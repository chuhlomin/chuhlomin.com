# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a static site generator for chuhlomin.com, written in Go. The generator processes Markdown files with YAML frontmatter, Go templates, and various media files to produce a static website optimized for CloudFlare Pages deployment.

## Build and Development Commands

### Build and Run Generator

```bash
# Build the generator (creates bin/generator)
make build-generator

# Run the generator (generates output/ directory)
make run-generator

# Build and run in one command
make build

# Clean all generated files
make clean
```

### Local Development

```bash
# Start MeiliSearch for local search functionality
docker compose up search -d

# Start Caddy web server (in separate terminal)
caddy run

# Build static site and update search index
make build
```

Open https://local.chuhlomin.com to view the site locally.

### Photo Processing

```bash
# Build and run photographer tool (processes photos, uploads to R2)
make photographer arguments="--titles-file=/path/to/titles.yml"
```

### Testing

```bash
# Run tests from generator directory
cd generator && go test
```

## Architecture

### Core Components

1. **generator/** - Main static site generator (Go application)
   - `main.go` - Entry point, configuration, and initialization
   - `generator.go` - Core generation logic with concurrent file processing
   - `markdown.go` - Markdown parsing and processing with YAML frontmatter
   - `templates.go` - Template helper functions
   - `og.go` - OpenGraph metadata fetching with caching
   - `photos.go` - Photo processing and sprite generation
   - `wishlist.go` - Wishlist YAML processing

2. **cmd/** - Command-line applications
   - **cmd/photographer/** - Photo upload and processing tool
     - Uploads photos to Cloudflare R2
     - Generates thumbnail sprites and blurhash placeholders
     - Updates `content/photos.yml`
   - **cmd/wishlist/** - Wishlist management tool

3. **content/** - Source content
   - Markdown files (`.md`) with YAML frontmatter
   - Go templates (`.gohtml`, `.goxml`, `.gotxt`)
   - YAML data files (`.yml`)
   - Static assets (CSS, images, fonts)

4. **templates/** - Reusable Go template partials
   - Used by content templates via template composition

5. **output/** - Generated static site (created by generator)

### File Processing Pipeline

The generator processes files concurrently using worker pools:

1. **File Discovery**: Walks `content/` directory to find all files
2. **Concurrent Processing** (4 workers by default):
   - `.md` files → Parse frontmatter, extract images, convert to HTML
   - `.gohtml`/`.goxml`/`.gotxt` → Process as Go templates
   - `.yml` files → Special handling (photos, wishlists)
   - Other files → Copy to output directory
3. **Image Processing**: Concurrent thumbnail generation with caching
4. **Template Rendering**: Render all Go templates with full site data
5. **Search Indexing**: Update MeiliSearch index by language

### Multi-language Support

Files are language-aware using underscore notation:
- `post.md` → English (default)
- `post_ru.md` → Russian
- Both files share the same `ID` (used for alternates)

The `ID` is extracted from the filename (everything before the language suffix). Templates can access alternate language versions via the `alternates` function.

### Template System

Go templates have access to:
- `Data` struct with:
  - `File` - Current MarkdownFile
  - `All` - Map of all MarkdownFiles by source path
  - `AllSorted` - All MarkdownFiles sorted by date
  - `Alternates` - Language alternatives for current page
  - `Timestamp` - Build timestamp

Available template functions (defined in `templates.go`):
- `config`, `enabled` - Access configuration values
- `i18n` - Translations from `i18n/` directory
- `date`, `link`, `alternates` - Formatting and navigation
- `md` - Render markdown to HTML inline

### Markdown Processing

Markdown files support:
- YAML frontmatter (date, title, tags, language, draft, etc.)
- Title extracted from first `# Header`
- Tags parsed from `#tag1 #tag2` lines
- Cross-references via `refs:` frontmatter
- Image extraction (both markdown and HTML syntax)
- Links to other `.md` files are rewritten to `.html` (or without extension if `REMOVE_HTML_EXT=true`)

### Configuration

Configuration is handled via environment variables or command-line flags (see `main.go` Config struct):
- `CONTENT_DIR` - Source content directory (default: `content`)
- `OUTPUT_DIR` - Generated files directory (default: `output`)
- `ROOT_URL` - Base URL (default: `https://local.chuhlomin.com`)
- `SEARCH_ENABLED` - Enable MeiliSearch integration
- `REMOVE_HTML_EXT` - Remove `.html` extension from URLs
- `SHOW_DRAFTS` - Include draft posts
- Many more in `main.go` Config struct

### CloudFlare Pages Specifics

The site is designed for CloudFlare Pages deployment:

1. **HTML Extension Removal**: `REMOVE_HTML_EXT=true` removes `.html` from links because CloudFlare Pages automatically redirects `/page.html` → `/page`

2. **Transform Rules** handle language query parameters:
   - `?lang=ru` rewrites to `_ru` suffix
   - Example: `/blog?lang=ru` → `/blog/index_ru.html`

3. **Redirects** in `_redirects` file handle legacy URLs and shortcuts

### Search Integration

When `SEARCH_ENABLED=true`:
- Generator indexes all MarkdownFiles to MeiliSearch
- Indexes are created per-language
- Uses `IDHash` as document identifier (alphanumeric-safe hash of ID)
- Frontend uses `SEARCH_API_KEY` for client-side search

### Caching Strategy

1. **Image Cache** (`cache/` directory):
   - Stores generated thumbnails
   - Checked before regenerating

2. **OpenGraph Cache** (`cache.yml`):
   - Caches fetched OpenGraph metadata
   - Persisted across builds

## Important Patterns

### Adding New Template Functions

Add to `fm` variable in `templates.go` and implement the function below.

### Adding New YAML Data Types

Update `processYaml()` in `generator.go` with a new case handler.

### Working with Photos

Photos are managed separately via the `photographer` tool, which:
1. Uploads full-res images to Cloudflare R2
2. Generates thumbnail sprites
3. Creates blurhash placeholders
4. Updates `content/photos.yml` with metadata

The generator consumes `photos.yml` to create photo gallery pages.

### Cross-references and Related Posts

Use `refs: [path/to/other.md]` in frontmatter to create relationships. The generator automatically creates bidirectional links.
