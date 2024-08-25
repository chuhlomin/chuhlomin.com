package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"gopkg.in/yaml.v3"
)

var (
	imageMarkdown  = regexp.MustCompile(`!\[(.*?)\]\(([^\s)]*)\s*"?([^"]*?)?"?\)`)
	imageHTML      = regexp.MustCompile(`<img(.*?)>`)
	htmlAttributes = regexp.MustCompile(`(\S+)\s*=\s*\"?(.*?)\"`)
	linkToMD       = regexp.MustCompile(`\[(.+?)\]\((.+?).md\)`)
)

// Data struct is used to pass data to the template
type Data struct {
	File       *MarkdownFile
	All        map[string]*MarkdownFile
	Timestamp  string
	AllSorted  []*MarkdownFile
	Alternates []*MarkdownFile // used only for index.html
}

// MarkdownFile represents a markdown file, for example
//
//	---
//	date: "2021-10-24"
//	---
//	# Title
//	Page content
type MarkdownFile struct {
	CommentsEnabled *bool    `yaml:"comments_enabled"`           // comments_enabled overrides config.CommentsEnabled
	Source          string   `yaml:"-"`                          // path to the source markdown file
	Path            string   `yaml:"-"`                          // path to the generated HTML file
	Canonical       string   `yaml:"-"`                          // canonical URL
	ID              string   `yaml:"-"`                          // same post in different languages will have the same ID value
	IDHash          string   `yaml:"-"`                          // hash of the ID, used as search index because: document identifier can be of type integer or string, only composed of alphanumeric characters (a-z A-Z 0-9), hyphens (-) and underscores (_)
	Markdown        string   `yaml:"-" indexer:"text"`           // content of the markdown file
	Title           string   `yaml:"title" indexer:"text"`       // by default equals to H1 in Markdown file
	Body            string   `yaml:"-" indexer:"no_store"`       // html body, generated from markdown
	Date            string   `yaml:"date" indexer:"date"`        // date when post was published, in format "2006-01-02"
	Type            string   `yaml:"type"`                       // "post" (by default), "page", etc.
	Language        string   `yaml:"language"`                   // language ("en", "ru", ...), parsed from filename, overrides config.DefaultLanguage
	Template        string   `yaml:"template"`                   // template to use in config.TemplatesDirectory, overrides default "post.html"
	Order           string   `yaml:"order"`                      // can be used to sort pages
	Description     string   `yaml:"description" indexer:"text"` // description is used for the meta description
	Author          string   `yaml:"author"`                     // author is used for the meta author
	Keywords        string   `yaml:"keywords"`                   // keywords is used for the meta keywords
	Image           string   `yaml:"image"`                      // image associated with the post; it's used to generate the thumbnailPath
	Tags            tags     `yaml:"tags"`                       // post tags, by default parsed from the post
	Refs            []string `yaml:"refs"`                       // references to other posts, used to generate the list of related posts
	Images          []image  `yaml:"-"`                          // images in the post
	Draft           bool     `yaml:"draft"`                      // draft is used to mark post as draft
}

// image is a struct that contains metadata of image from the post
type image struct {
	Path      string `yaml:"path"`
	Alt       string `yaml:"alt"`
	Title     string `yaml:"title"`
	ThumbPath string `yaml:"thumb_path"`
	Promo     bool   `yaml:"promo"`
}

type tags []string

func (t *tags) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	*t = strings.Split(s, ", ")
	return nil
}

type ByCreated []*MarkdownFile

func (md ByCreated) Len() int           { return len(md) }
func (md ByCreated) Less(i, j int) bool { return md[i].Date > md[j].Date }
func (md ByCreated) Swap(i, j int)      { md[i], md[j] = md[j], md[i] }

type ByOrder []*MarkdownFile

func (md ByOrder) Len() int           { return len(md) }
func (md ByOrder) Less(i, j int) bool { return md[i].Order < md[j].Order }
func (md ByOrder) Swap(i, j int)      { md[i], md[j] = md[j], md[i] }

type ByLanguage []*MarkdownFile

func (md ByLanguage) Len() int           { return len(md) }
func (md ByLanguage) Less(i, j int) bool { return md[i].Language > md[j].Language }
func (md ByLanguage) Swap(i, j int)      { md[i], md[j] = md[j], md[i] }

func NewMarkdownFile(dir, path string) (*MarkdownFile, error) {
	content, err := os.ReadFile(filepath.Join(dir, path))
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return processMarkdownFileContent(path, content)
}

func processMarkdownFileContent(path string, content []byte) (*MarkdownFile, error) {
	outputPath := strings.Replace(path, ".md", ".html", 1)
	outputPath = strings.TrimPrefix(outputPath, cfg.ContentDirectory)

	md := &MarkdownFile{
		Source:    path,
		Path:      "/" + outputPath,
		Canonical: cfg.RootURL + "/" + link(outputPath),
		Tags:      tags([]string{}), // setting default value, so that there is no need to check for nil in templates
	}

	md.ID, md.Language = getIDAndLangFromPath(path)
	md.IDHash = fmt.Sprintf("%x", sha256.Sum256([]byte(md.ID)))
	if md.Language == "" {
		md.Language = cfg.DefaultLanguage
	}
	md.Language = strings.ToLower(md.Language)

	if md.CommentsEnabled == nil {
		md.CommentsEnabled = &cfg.CommentsEnabled
	}

	bodyBytes, err := md.processContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to process content of %s: %w", path, err)
	}

	md.Markdown = string(bodyBytes)

	bodyBytes = markdown.ToHTML(bodyBytes, nil, nil)

	// todo: add typograph here :typograph.NewTypograph().Process(bodyBytes)

	md.Body = string(bodyBytes)

	return md, nil
}

func (md *MarkdownFile) processContent(content []byte) ([]byte, error) {
	metadataBytes, bodyBytes := splitMetadataAndBody(content)

	if err := md.parseMetadata(metadataBytes); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	baseDir := filepath.Dir(md.Source)
	relativePath := baseDir
	thumbPath := cfg.ThumbPath + "/" + baseDir

	if md.Image != "" {
		// Add image to the list of all images
		path, thumbPath := fixPath(md.Image, relativePath, thumbPath)
		md.Images = append(md.Images, image{
			Path:      path,
			ThumbPath: thumbPath,
		})
	}

	bodyBytes = md.processBody(bodyBytes, relativePath, thumbPath)

	return bodyBytes, nil
}

// processBody parses Title, Tags and Images from Markdown file content.
// It removes Title and Tags from the content,
// so that they are not processed by markdown.ToHTML call.
// relativePath and thumbPath are needed for image processing and subject to refactoring.
func (md *MarkdownFile) processBody(b []byte, relativePath, thumbPath string) []byte {
	b = bytes.TrimSpace(b)

	buf := bytes.Buffer{}
	hasHeader := false

	var tags []string

	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		line := scanner.Text()
		b := scanner.Bytes()

		// parse header
		if strings.HasPrefix(line, "# ") && !hasHeader {
			htmlTitle := string(markdown.ToHTML([]byte(strings.TrimSpace(line[2:])), nil, nil))
			htmlTitle = strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(htmlTitle), "<p>"), "</p>")

			md.Title = htmlTitle
			hasHeader = true
			continue // so that we don't leave header in the body
		}

		// parse tags
		if len(md.Tags) == 0 && strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "# ") && !strings.HasPrefix(line, "##") {
			tags = strings.Split(strings.TrimSpace(line), " ")
			for i, tag := range tags {
				tags[i] = strings.Trim(tag, "#,")
			}
			md.Tags = tags
			continue // so that we don't leave tags in the body
		}

		// parse markdown images
		if matches := imageMarkdown.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				alt, url, title := match[1], match[2], match[3]

				path, thumbPath := fixPath(url, relativePath, thumbPath)
				md.Images = append(md.Images, image{
					Path:      path,
					Alt:       alt,
					Title:     title,
					ThumbPath: thumbPath,
				})
			}
		}

		// parse HTML images
		if matches := imageHTML.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				img := image{}
				attributes := htmlAttributes.FindAllStringSubmatch(match[1], -1)
				for _, attr := range attributes {
					switch attr[1] {
					case "src":
						path, thumbPath := fixPath(attr[2], relativePath, thumbPath)
						img.Path = path
						img.ThumbPath = thumbPath
					case "alt":
						img.Alt = attr[2]
					case "title":
						img.Title = attr[2]
					}
				}

				md.Images = append(md.Images, img)
			}
		}

		// parse links to markdown files
		// replace links to markdown files with links to HTML files
		if matches := linkToMD.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				path := link(match[2]+".html", md.Language)
				line = strings.Replace(line, match[0], "["+match[1]+"]("+path+")", -1)
				b = []byte(line)
			}
		}

		buf.Write(b)
		buf.WriteString("\n")
	}

	return buf.Bytes()
}

// parseMetadata parses YAML metadata at the beginning of the Markdown file.
// For convinience, it puts all metadata into MarkdownFile struct.
func (md *MarkdownFile) parseMetadata(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	if err := yaml.Unmarshal(b, &md); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return nil
}

func splitMetadataAndBody(b []byte) ([]byte, []byte) {
	if bytes.HasPrefix(b, []byte("---")) {
		if parts := bytes.SplitN(b, []byte("---"), 3); len(parts) == 3 {
			return parts[1], parts[2]
		}
	}

	return []byte{}, b
}

func getIDAndLangFromPath(filename string) (id, lang string) {
	underscoreIndex := strings.LastIndex(filename, "_")
	if underscoreIndex == -1 {
		return filename, ""
	}

	dotIndex := strings.LastIndex(filename, ".")
	if dotIndex == -1 {
		return filename, ""
	}

	lang = filename[underscoreIndex+1 : dotIndex]
	if len(lang) != 2 {
		return filename, ""
	}

	id = filename[0:underscoreIndex] + filename[dotIndex:]
	return
}

func fixPath(url, relativePath, thumbPath string) (string, string) {
	if !isValidURL(url) {
		return path.Clean(relativePath + "/" + url),
			path.Clean(thumbPath + "/" + url)
	}
	// sha1 hash of the url
	h := sha1.New()
	h.Write([]byte(url))
	hash := hex.EncodeToString(h.Sum(nil))

	// get path file extension
	ext := filepath.Ext(url)

	return url, thumbPath + "/" + hash + ext
}

func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
