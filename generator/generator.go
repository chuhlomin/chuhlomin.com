package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	goimage "image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/charmbracelet/log"
	"github.com/disintegration/imaging"
	"github.com/meilisearch/meilisearch-go"
)

const (
	permDir  = 0755 // permission used to create directories in cfg.OutputDirectory
	permFile = 0644 // permissions for
)

// Generator holds all the data needed to generate the site
type Generator struct {
	cfg          Config
	t            *template.Template
	md           map[string]*MarkdownFile // relative source path -> MarkdownFile
	mdSorted     []*MarkdownFile          // sorted MarkdownFiles by creation date
	mdMu         sync.Mutex
	tempDir      string                              // temporary directory used to templates
	templates    map[string]map[string]*MarkdownFile // id -> hashed path -> MarkdownFile
	templatesMu  sync.Mutex
	searchClient *meilisearch.Client
	og           *openGraphClient
}

// NewGenerator creates a new Generator
func NewGenerator(
	ogClient *openGraphClient,
	searchClient *meilisearch.Client,
) (*Generator, error) {
	t, err := template.New("").Funcs(fm).ParseGlob(cfg.TemplatesDirectory + "/*")
	if err != nil {
		return nil, fmt.Errorf("Error parsing templates: %v", err)
	}

	tempDir, err := os.MkdirTemp(cfg.TempDirectory, "templates")
	if err != nil {
		return nil, fmt.Errorf("Error creating temp directory: %v", err)
	}

	return &Generator{
		t:            t,
		md:           map[string]*MarkdownFile{},
		tempDir:      tempDir,
		templates:    map[string]map[string]*MarkdownFile{},
		templatesMu:  sync.Mutex{},
		searchClient: searchClient,
		og:           ogClient,
	}, nil
}

// Run runs the generator, generating the site
func (g *Generator) Run(ts time.Time) error {
	var (
		files              = make(chan string, cfg.FilesChannelSize)
		images             = make(chan image, cfg.ImagesChannelSize)
		doneFiles          = make(chan bool)
		doneImages         = make(chan bool)
		doneI18s           = make(chan bool)
		doneSearchIndexing = make(chan bool)
	)

	go walkDir(cfg.ContentDirectory, files)
	go g.processFiles(files, images, doneFiles)
	go g.processI18n(doneI18s)
	go g.processImages(images, doneImages)

	<-doneFiles

	// add cross references
	for _, file := range g.md {
		if file.Refs == nil {
			continue
		}

		for _, ref := range file.Refs {
			refFile, ok := g.md[ref]
			if !ok {
				continue
			}

			if !contains(refFile.Refs, file.Source) {
				refFile.Refs = append(refFile.Refs, file.Source)
			}
		}
	}

	// sort MarkdownFiles by creation date
	g.mdSorted = make([]*MarkdownFile, 0, len(g.md))
	for _, file := range g.md {
		g.mdSorted = append(g.mdSorted, file)
	}
	sort.Sort(ByCreated(g.mdSorted))

	if g.isSearchEnabled() {
		go g.updateSearchIndex(doneSearchIndexing)
	}

	<-doneI18s
	g.renderAllTemplates(ts)
	g.renderAllMarkdown(ts)

	log.Info("Waiting for images to be processed...")
	<-doneImages

	if g.isSearchEnabled() {
		log.Info("Waiting for search index to be updated...")
		<-doneSearchIndexing
	}

	return nil
}

func (g *Generator) processFiles(files <-chan string, images chan<- image, done chan<- bool) {
	wg := sync.WaitGroup{}

	for i := 0; i < cfg.NumWorkers; i++ {

		for path := range files {
			go func(path string) {
				wg.Add(1)
				defer wg.Done()

				if err := g.processFile(path, images); err != nil {
					log.Fatalf("Error processing file %s: %v", path, err)
				}
			}(path)
		}
	}

	wg.Wait()

	done <- true
	close(images)
}

func (g *Generator) processImages(images <-chan image, done chan<- bool) {
	processedImages := make(map[string]bool)

	wg := sync.WaitGroup{}
	accessMap := sync.Mutex{}

	for i := 0; i < cfg.NumWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for img := range images {
				accessMap.Lock()
				if processedImages[img.Path] {
					accessMap.Unlock()
					continue
				}
				accessMap.Unlock()

				if err := g.processImage(img); err != nil {
					log.Fatalf("Error processing image %s: %v", img.Path, err)
				}

				accessMap.Lock()
				processedImages[img.Path] = true
				accessMap.Unlock()
			}
		}()
	}

	wg.Wait()
	done <- true
}

func (g *Generator) processI18n(done chan<- bool) {
	filepath.Walk(cfg.I18NDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Optional .toml files are used to define translations.
		// They power `i18n` template function.
		_, err = bundle.LoadMessageFile(path)
		if err != nil {
			log.Fatalf("ERROR load message file %q: %v", path, err)
		}
		return nil
	})

	done <- true
}

func (g *Generator) processFile(path string, images chan<- image) error {
	ext := filepath.Ext(path)

	switch ext {
	case ".md":
		return g.processMarkdown(path, images)
	case ".yml":
		return g.processYaml(path)
	default:
		if strings.HasPrefix(ext, ".go") {
			return g.processGoTemplate(path)
		}

		if err := g.copyFile(
			filepath.Join(cfg.ContentDirectory, path),
			filepath.Join(cfg.OutputDirectory, path),
		); err != nil {
			return fmt.Errorf("Error copying file %s: %v", path, err)
		}
		return nil
	}
}

func (g *Generator) processGoTemplate(path string) error {
	log.Debugf("Processing Go template %s", path)

	// use sha256 hash the path to the template
	// hashing is used to avoid collisions in the template package
	// From the docs:
	//   When parsing multiple files with the same name in different directories,
	//   the last one mentioned will be the one that results.
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(path)))

	id, lang := getIDAndLangFromPath(path)
	if lang == "" {
		lang = cfg.DefaultLanguage
	}

	ext := filepath.Ext(path) // .gohtml, .goxml, .gotxt, etc.
	// newExt is the extension of the output file
	// .html, .xml, .txt, etc.
	newExt := "." + ext[3:]
	outputPath := filepath.Join(
		"/",
		filepath.Dir(path),
		filepath.Base(path[:len(path)-len(ext)])+newExt,
	)

	g.templatesMu.Lock()
	if g.templates[id] == nil {
		g.templates[id] = map[string]*MarkdownFile{}
	}
	g.templates[id][hash] = &MarkdownFile{
		Source:   path,
		ID:       id,
		IDHash:   fmt.Sprintf("%x", sha256.Sum256([]byte(id))),
		Path:     outputPath,
		Language: lang,
	}
	g.templatesMu.Unlock()

	// copy the template to the temp directory
	if err := g.copyFile(
		filepath.Join(cfg.ContentDirectory, path),
		filepath.Join(g.tempDir, hash),
	); err != nil {
		return fmt.Errorf("Error copying template %s: %v", path, err)
	}

	return nil
}

func (g *Generator) processMarkdown(path string, images chan<- image) error {
	log.Debugf("Processing Markdown %s", path)

	// Markdown files are posts/pages and are rendered as HTML using templates
	// All posts are stored in g.md and then rendered, so that we can use the
	// file metadata to generate the index pages, etc.

	md, err := NewMarkdownFile(cfg.ContentDirectory, path)
	if err != nil {
		return fmt.Errorf("Error creating MarkdownFile: %v", err)
	}

	if md.Draft && !cfg.ShowDrafts {
		log.Debugf("Skipping draft %v", path)
		return nil
	}

	for _, image := range md.Images {
		images <- image
	}

	g.mdMu.Lock()
	g.md[md.Source] = md
	g.mdMu.Unlock()

	return nil
}

func (g *Generator) processYaml(path string) error {
	log.Debugf("Processing Yaml %s", path)

	fileContent, err := os.ReadFile(filepath.Join(cfg.ContentDirectory, path))
	if err != nil {
		return fmt.Errorf("Error reading file %s: %v", path, err)
	}

	outputPath := filepath.Join(cfg.OutputDirectory, strings.Replace(path, ".yml", ".html", 1))

	templateName := strings.Replace(path, ".yml", ".gohtml", 1)
	tmpl := g.t.Lookup(templateName)
	if tmpl == nil {
		return fmt.Errorf("Template %q not found", templateName)
	}

	var data interface{}

	switch path {
	case "wishlist.yml":
		data, err = g.processWishlistItems(fileContent)
		if err != nil {
			return fmt.Errorf("Error processing wishlist items for %q: %v", path, err)
		}
	case "photos.yml":
		return nil // noop

		// data, err = g.processPhotos(fileContent)
		// if err != nil {
		// 	return fmt.Errorf("Error processing photos for %q: %v", path, err)
		// }
	default:
		log.Fatalf("Unknown YAML file %q", path)
	}

	return g.renderTemplate(outputPath, data, tmpl)
}

func (g *Generator) processImage(image image) error {
	log.Debugf("Processing image %s", image.Path)

	if strings.HasSuffix(image.Path, ".svg") {
		// SVGs are copied as-is
		if err := g.copyFile(
			filepath.Join(cfg.ContentDirectory, image.Path),
			filepath.Join(cfg.OutputDirectory, image.Path),
		); err != nil {
			return fmt.Errorf("Error copying file %s: %v", image.Path, err)
		}
		return nil
	}

	// check cache directory first
	if _, err := os.Stat(filepath.Join(cfg.CacheDirectory, image.ThumbPath)); err == nil {
		log.Debugf("Image %s already exists in cache", image.ThumbPath)
		err = g.copyFile(
			filepath.Join(cfg.CacheDirectory, image.ThumbPath),
			filepath.Join(cfg.OutputDirectory, image.ThumbPath),
		)
		if err != nil {
			return fmt.Errorf("Error copying image %s: %v", image.Path, err)
		}
		return nil
	}

	var (
		img goimage.Image
		err error
	)

	// read image
	if isValidURL(image.Path) {
		img, err = getImageFromURL(image.Path)
		if err != nil {
			return fmt.Errorf("get image from url %q: %v", image.Path, err)
		}
	} else {
		img, err = imaging.Open(
			filepath.Join(cfg.ContentDirectory, image.Path),
			imaging.AutoOrientation(true),
		)
		if err != nil {
			return fmt.Errorf("read image %q: %v", image.Path, err)
		}
	}

	// resize image
	img = imaging.Fit(img, cfg.ThumbMaxWidth, cfg.ThumbMaxHeight, imaging.Lanczos)

	// save image
	thumbPath := filepath.Join(cfg.OutputDirectory, image.ThumbPath)
	if err = os.MkdirAll(filepath.Dir(thumbPath), permDir); err != nil {
		return fmt.Errorf("create directory %q: %v", filepath.Dir(thumbPath), err)
	}
	if err := imaging.Save(img, thumbPath); err != nil {
		return fmt.Errorf("save image %q: %v", thumbPath, err)
	}

	// copy image to cache
	err = g.copyFile(
		filepath.Join(cfg.OutputDirectory, image.ThumbPath),
		filepath.Join(cfg.CacheDirectory, image.ThumbPath),
	)
	if err != nil {
		return fmt.Errorf("Error copying image to cache %s: %v", image.ThumbPath, err)
	}

	return nil
}

func (g *Generator) renderAllTemplates(ts time.Time) {
	var err error
	g.t, err = g.t.ParseGlob(g.tempDir + "/*")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// iterate over g.templates
	for _, hashFile := range g.templates {
		alternates := []*MarkdownFile{}

		for _, file := range hashFile {
			alternates = append(alternates, file)
		}
		sort.Sort(ByLanguage(alternates))

		for hash, file := range hashFile {
			err := g.renderGoTemplate(hash, file, alternates, ts)
			if err != nil {
				log.Fatalf("Error rendering template %s: %v", file.Path, err)
			}
		}
	}
}

func (g *Generator) renderGoTemplate(
	hash string,
	file *MarkdownFile,
	alternates []*MarkdownFile,
	ts time.Time,
) error {
	log.Debugf("Rendering %s → %s", file.Source, file.Path)
	t := g.t.Lookup(hash)
	if t == nil {
		return fmt.Errorf("Template %s not found", file.Source)
	}

	return g.renderTemplate(
		filepath.Join(cfg.OutputDirectory, file.Path),
		Data{
			File:       file,
			Alternates: alternates,
			All:        g.md,
			AllSorted:  g.mdSorted,
			Timestamp:  ts.Format(time.RFC3339),
		},
		t,
	)
}

func (g *Generator) renderAllMarkdown(ts time.Time) {
	wg := sync.WaitGroup{}

	for i := 0; i < cfg.NumWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for _, md := range g.md {
				err := g.renderMarkdown(md, ts)
				if err != nil {
					log.Fatalf("Error rendering Markdown %s: %v", md.Path, err)
				}
			}
		}()
	}

	wg.Wait()
}

func (g *Generator) copyFile(src, dst string) error {
	log.Debugf("Copying %s → %s", src, dst)

	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, permDir); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

func (g *Generator) renderMarkdown(md *MarkdownFile, ts time.Time) error {
	tmpl := g.t.Lookup("post.gohtml")
	// if file.Template != "" {
	// 	tmpl = defaultTmpl.Lookup(file.Template)
	// 	if tmpl == nil {
	// 		return errors.Errorf("template %q not found", file.Template)
	// 	}
	// }

	if err := g.renderTemplate(
		filepath.Join(cfg.OutputDirectory, md.Path),
		Data{
			File:      md,
			All:       g.md,
			AllSorted: g.mdSorted,
			Timestamp: ts.Format(time.RFC3339),
		},
		tmpl,
	); err != nil {
		return fmt.Errorf("Error rendering template: %v", err)
	}

	return nil
}

func (g *Generator) renderTemplate(outputPath string, data interface{}, t *template.Template) error {
	// create directories for file
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll("./"+dir, permDir); err != nil {
		return fmt.Errorf("creating directory %s: %v", dir, err)
	}

	// open file
	f, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, permFile)
	if err != nil {
		return fmt.Errorf("opening file %s: %v", outputPath, err)
	}
	defer f.Close()

	// execute template
	if err = t.Execute(f, data); err != nil {
		return fmt.Errorf("executing template: %v", err)
	}

	return nil
}

func (g *Generator) updateSearchIndex(doneSearchIndexing chan<- bool) {
	languageBatches := map[string][]*MarkdownFile{}
	for _, file := range g.md {
		languageBatches[file.Language] = append(languageBatches[file.Language], file)
	}

	for language, files := range languageBatches {
		tasks, err := g.searchClient.Index(language).AddDocumentsInBatches(files, 1000, "IDHash")
		if err != nil {
			log.Fatalf("Error indexing documents: %v", err)
		}

		for _, task := range tasks {
			// increase default context timeout from 5s to 2m to wait for slow
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute*2)
			defer cancelFunc()
			_, err := g.searchClient.WaitForTask(
				task.TaskUID,
				meilisearch.WaitParams{
					Context:  ctx,
					Interval: time.Millisecond * 500,
				},
			)
			if err != nil {
				log.Fatalf("Error waiting for task %d: %v", task.TaskUID, err)
			}
		}
	}

	doneSearchIndexing <- true
}

func (g *Generator) isSearchEnabled() bool {
	return g.searchClient != nil
}

func walkDir(dir string, files chan<- string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		files <- path[len(dir)+1:]
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory %s: %v", dir, err)
	}

	close(files)
}

func getImageFromURL(url string) (goimage.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get image from url")
	}
	defer resp.Body.Close()

	img, _, err := goimage.Decode(resp.Body)
	return img, err
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
