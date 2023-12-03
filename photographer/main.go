// A small go app to handle photos
// It will
// - generate thumbnails for all photos
// - upload photos and thumbnails to R2 Cloudflare storage if they don't exist yet
// - generate BlurHash for each photo
// - keep track of them in a photos.yaml file that will be used by the Generator to generate the website
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"image"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bbrks/go-blurhash"
	"github.com/charmbracelet/log"
	flags "github.com/jessevdk/go-flags"
	"github.com/nfnt/resize"
	"gopkg.in/yaml.v3"

	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
)

// Photo struct for items in photos.yml file
// Must be in sync with generator/photos.go
type Photo struct {
	Path                string
	Title               string `yaml:"title,omitempty"`
	TitleRu             string `yaml:"title_ru,omitempty"`
	Width               int    `yaml:"width,omitempty"`
	Height              int    `yaml:"height,omitempty"`
	ThumbPath           string `yaml:"thumb,omitempty"`
	ThumbXOffset        int    `yaml:"thumb_x,omitempty"`
	ThumbYOffset        int    `yaml:"thumb_y,omitempty"`
	ThumbWidth          int    `yaml:"thumb_width,omitempty"`
	ThumbHeight         int    `yaml:"thumb_height,omitempty"`
	ThumbTotalWidth     int    `yaml:"thumb_total_width,omitempty"`
	ThumbTotalHeight    int    `yaml:"thumb_total_height,omitempty"`
	Blurhash            string `yaml:"blurhash,omitempty"`
	BlurhashImageBase64 string `yaml:"blurhash_image_base64,omitempty"`

	// Temporary image.Image field used to generate thumbnails
	image image.Image `yaml:"-"`
}

// PhotoContainer is a wrapper for Photo struct, used for sorting,
// so that references are not swapped and still can be modified
type PhotoContainer struct {
	Photo *Photo
}

type byThumbHeightDesc []PhotoContainer

func (a byThumbHeightDesc) Len() int      { return len(a) }
func (a byThumbHeightDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byThumbHeightDesc) Less(i, j int) bool {
	return a[i].Photo.ThumbHeight > a[j].Photo.ThumbHeight
}

type appConfig struct {
	// App uses `photo.yml` file as a state file to keep track of photos
	YamlFile string `env:"YAML_FILE" long:"yaml-file" description:"path to photos.yml file" default:"content/photos.yml"`

	// Directory with photos
	PhotosDir string `env:"PHOTOS_DIR" long:"photos-dir" description:"path to photos directory" default:"~/Pictures/Photos/"`

	// Cloudflare R2 storage
	R2AccountID       string `env:"R2_ACCOUNT_ID" long:"r2-account-id" description:"r2 account id"`
	R2AccessKeyID     string `env:"R2_ACCESS_KEY_ID" long:"r2-access-key-id" description:"r2 access key id"`
	R2AccessKeySecret string `env:"R2_ACCESS_KEY_SECRET" long:"r2-access-key-secret" description:"r2 access key secret"`
	R2Bucket          string `env:"R2_BUCKET" long:"r2-bucket" description:"r2 bucket"`

	// Titles file from Instagram, to populate titles for photos
	TitlesFile string `env:"TITLES_FILE" long:"titles-file" description:"path to posts_1.json file"`

	// Force thumbnail generation
	ForceThumbnails bool `long:"force-thumbnails" description:"force thumbnail generation"`

	// Blurhash
	ForceBlurhash       bool `long:"force-blurhash" description:"force blurhash generation"`
	ForceBlurhashImages bool `long:"force-blurhash-images" description:"force blurhash images generation"`
}

func main() {
	log.Info("Starting...")

	if err := run(); err != nil {
		log.Fatal(err)
	}

	log.Info("Finished")
}

func run() error {
	var cfg appConfig
	_, err := flags.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("error parsing flags: %v", err)
	}

	ctx := context.Background()

	photos, err := loadPhotosFile(cfg.YamlFile)
	if err != nil {
		return fmt.Errorf("error loading photos: %v", err)
	}

	dir, err := absolutePath(cfg.PhotosDir)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}

	files, err := scanDirectory(dir)
	if err != nil {
		return fmt.Errorf("error scanning directory: %v", err)
	}

	r2, err := NewR2(
		cfg.R2AccountID,
		cfg.R2AccessKeyID,
		cfg.R2AccessKeySecret,
		cfg.R2Bucket,
	)
	if err != nil {
		return fmt.Errorf("error creating r2 client: %v", err)
	}

	photos, err = uploadNewPhotos(ctx, r2, photos, files, dir)
	if err != nil {
		return fmt.Errorf("error uploading new photos: %v", err)
	}

	photos, err = generateThumbnails(ctx, r2, photos, dir, cfg.ForceThumbnails)
	if err != nil {
		return fmt.Errorf("error generating thumbnails: %v", err)
	}

	photos, err = generateBlurhashes(photos, dir, cfg.ForceBlurhash)
	if err != nil {
		return fmt.Errorf("error generating blurhashes: %v", err)
	}

	photos, err = generateBlurhashImages(photos, cfg.ForceBlurhashImages)
	if err != nil {
		return fmt.Errorf("error generating blurhash images: %v", err)
	}

	if cfg.TitlesFile != "" {
		photos, err = populateTitles(photos, cfg.TitlesFile)
		if err != nil {
			return fmt.Errorf("error populating titles: %v", err)
		}
	}

	// save photos.yml file
	if err = savePhotosFile(cfg.YamlFile, photos); err != nil {
		return fmt.Errorf("error saving photos: %v", err)
	}

	return nil
}

func loadPhotosFile(path string) ([]*Photo, error) {
	// read photos.yml file
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var photos []*Photo
	if err = yaml.Unmarshal(fileContent, &photos); err != nil {
		return nil, fmt.Errorf("error unmarshaling file: %v", err)
	}

	return photos, nil
}

func absolutePath(dir string) (string, error) {
	// directory might be relative to home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting absolute path: %v", err)
	}

	// replace ~ with home directory
	return strings.Replace(dir, "~", home, 1), nil
}

func scanDirectory(dir string) ([]string, error) {
	var files []string
	// read directory recursively
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Skip files that are JPG
		if filepath.Ext(path) != ".jpg" && filepath.Ext(path) != ".jpeg" {
			return nil
		}

		// Ignore thumbnails (they are handled separately)
		if strings.HasPrefix(filepath.Base(path), "thumbnails_") {
			return nil
		}

		// trim the dir prefix
		path = strings.TrimPrefix(path, dir)
		files = append(files, path)
		return nil
	})

	return files, err
}

func diff(photos []*Photo, files []string) ([]string, []string) {
	var toAdd []string
	var toDelete []string

	// find new files
	for _, file := range files {
		if !containsPhoto(photos, file) {
			toAdd = append(toAdd, file)
		}
	}

	// find deleted files
	for _, photo := range photos {
		if !contains(files, photo.Path) {
			toDelete = append(toDelete, photo.Path)
		}
	}

	return toAdd, toDelete
}

func contains(arr []string, needle string) bool {
	for _, item := range arr {
		if item == needle {
			return true
		}
	}

	return false
}

func containsPhoto(arr []*Photo, needle string) bool {
	for _, item := range arr {
		if item.Path == needle {
			return true
		}
	}

	return false
}

func savePhotosFile(path string, photos []*Photo) error {
	// marshal photos to yaml
	fileContent, err := yaml.Marshal(photos)
	if err != nil {
		return fmt.Errorf("error marshaling photos: %v", err)
	}

	// save photos.yml file
	err = os.WriteFile(path, fileContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

func uploadNewPhotos(
	ctx context.Context,
	r2 *R2,
	photos []*Photo,
	files []string,
	dir string,
) ([]*Photo, error) {
	toAdd, toDelete := diff(photos, files)

	for _, file := range toAdd {
		photos = append(photos, &Photo{
			Path: file,
		})

		content, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		log.Infof("Uploading %s", file)
		if err = r2.Upload(ctx, file, content); err != nil {
			return nil, fmt.Errorf("error uploading file: %v", err)
		}
	}

	for _, file := range toDelete {
		for i, photo := range photos {
			if photo.Path == file {
				photos = append(photos[:i], photos[i+1:]...)
				break
			}
		}
	}

	return photos, nil
}

const (
	maxThumbSize = 324 /* 162 * 2 */
	maxPerRow    = 10
	maxRows      = 5
)

func generateThumbnails(
	ctx context.Context,
	r2 *R2,
	photos []*Photo,
	dir string,
	force bool,
) ([]*Photo, error) {
	// split photos into batches of 100 photos each
	photosBatches := make([][]*Photo, 0)
	for i := 0; i < len(photos); i += maxPerRow * maxRows {
		end := i + maxPerRow*maxRows
		if end > len(photos) {
			end = len(photos)
		}
		photosBatches = append(photosBatches, photos[i:end])
	}

	// filter out batches if all photos in it already have thumbnails
	if !force {
		for batch, photos := range photosBatches {
			allHaveThumbs := true
			allHaveSameThumb := true
			for _, photo := range photos {
				if photo.ThumbPath == "" {
					allHaveThumbs = false
					break
				}
				if photo.ThumbPath != photos[0].ThumbPath {
					allHaveSameThumb = false
					break
				}
			}
			if allHaveThumbs && allHaveSameThumb {
				photosBatches[batch] = nil
			}
		}
	}

	// generate thumbnails for each year
	for batch, photos := range photosBatches {
		if photos == nil {
			continue
		}

		thumbPath, err := generateThumbnail(batch, photos, dir)
		if err != nil {
			return nil, fmt.Errorf("error generating thumbnail for %d: %v", batch, err)
		}

		// upload thumbnail to R2
		thumbContent, err := os.ReadFile(filepath.Join(dir, thumbPath))
		if err != nil {
			return nil, fmt.Errorf("error reading thumbnail %q: %v", thumbPath, err)
		}

		if err := r2.Upload(ctx, thumbPath, thumbContent); err != nil {
			return nil, fmt.Errorf("error uploading thumbnail %q: %v", thumbPath, err)
		}

		// update thumb path with CRC32 checksum for each photo
		for _, photo := range photos {
			photo.ThumbPath = thumbPath + "?crc=" + crc32sum(thumbContent)
		}
	}

	return photos, nil
}

func generateThumbnail(batch int, photos []*Photo, dir string) (string, error) {
	log.Infof("Generating thumbnail for %d", batch)
	// each thumbnail should fit into 140x140px square, maximum 10 photos in a row
	for _, photo := range photos {
		// decode photo
		img, err := readImage(dir, photo.Path)
		if err != nil {
			return "", fmt.Errorf("error reading image: %v", err)
		}
		photo.Width = img.Bounds().Dx()
		photo.Height = img.Bounds().Dy()

		// resize photo to 140x140px
		img = resize.Thumbnail(
			maxThumbSize,
			maxThumbSize,
			img,
			resize.Lanczos3,
		)
		photo.image = img
		photo.ThumbWidth = img.Bounds().Dx()
		photo.ThumbHeight = img.Bounds().Dy()
	}

	// sort photos by height, aiming to have less empty space
	// create a slice of pointers to the original photos
	containers := make([]PhotoContainer, len(photos))
	for i := range photos {
		containers[i].Photo = photos[i]
	}

	// sort the slice of pointers by thumb height in descending order
	sort.Sort(byThumbHeightDesc(containers))

	// calculate thumbnail image size
	var (
		width   int
		height  int
		counter int
	)
	for i, container := range containers {
		if i == 0 {
			width = container.Photo.ThumbWidth
			height = container.Photo.ThumbHeight
		}

		if counter == maxPerRow {
			counter = 0
			height += container.Photo.ThumbHeight
		}

		if i < maxPerRow-1 {
			width += container.Photo.ThumbWidth
		}

		counter++
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// draw photos on thumbnail
	var (
		thumbPath = "thumbnails_" + strconv.Itoa(batch) + ".jpg"
		x         int
		y         int
		col       int
		rowHeight int
	)

	for i, container := range containers {
		if i == 0 {
			rowHeight = container.Photo.ThumbHeight
		}

		if col == maxPerRow {
			x = 0
			col = 0
			y += rowHeight
			rowHeight = container.Photo.ThumbHeight
		}

		container.Photo.ThumbPath = thumbPath
		container.Photo.ThumbXOffset = x
		container.Photo.ThumbYOffset = y
		container.Photo.ThumbTotalWidth = width
		container.Photo.ThumbTotalHeight = height

		draw.Draw(
			img,
			image.Rect(x, y, x+container.Photo.ThumbWidth, y+container.Photo.ThumbHeight),
			container.Photo.image,
			image.Point{0, 0},
			draw.Src,
		)
		x += container.Photo.ThumbWidth
		col++
	}

	// encode img thumbnail into JPEG
	out, err := os.Create(filepath.Join(dir, thumbPath))
	if err != nil {
		return "", fmt.Errorf("error creating file %q: %v", thumbPath, err)
	}
	defer out.Close()

	jpegOptions := jpeg.Options{
		Quality: 90,
	}
	if err = jpeg.Encode(out, img, &jpegOptions); err != nil {
		return "", fmt.Errorf("error encoding thumbnail: %v", err)
	}

	return thumbPath, nil
}

func readImage(dir string, path string) (image.Image, error) {
	file, err := os.Open(filepath.Join(dir, path))
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	return img, nil
}

func generateBlurhashes(photos []*Photo, dir string, force bool) ([]*Photo, error) {
	var err error
	for _, photo := range photos {
		if photo.Blurhash != "" && !force {
			continue
		}

		log.Infof("Generating blurhash for %s", photo.Path)
		photo.Blurhash, err = generateBlurhash(photo.Path, dir)
		if err != nil {
			return nil, fmt.Errorf("error generating blurhash: %v", err)
		}
	}

	return photos, nil
}

func generateBlurhash(path, dir string) (string, error) {
	file, err := os.Open(filepath.Join(dir, path))
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	return generateBlurhashForReader(file)
}

func generateBlurhashForReader(reader io.Reader) (string, error) {
	m, _, err := image.Decode(reader)
	if err != nil {
		return "", err
	}

	return blurhash.Encode(4, 4, m)
}

func generateBlurhashImages(photos []*Photo, force bool) ([]*Photo, error) {
	var err error
	for _, photo := range photos {
		if photo.Blurhash == "" {
			continue
		}

		if photo.BlurhashImageBase64 != "" && !force {
			continue
		}

		log.Infof("Generating blurhash image for %s", photo.Path)
		photo.BlurhashImageBase64, err = generateBlurhashImage(photo)
		if err != nil {
			return nil, fmt.Errorf("error generating blurhash image: %v", err)
		}
	}

	return photos, nil
}

func generateBlurhashImage(photo *Photo) (string, error) {
	m, err := blurhash.Decode(
		photo.Blurhash,
		photo.ThumbWidth/2,
		photo.ThumbHeight/2,
		1,
	)
	if err != nil {
		return "", fmt.Errorf("error decoding blurhash: %v", err)
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, m, &jpeg.Options{Quality: 90}); err != nil {
		return "", fmt.Errorf("error encoding blurhash image: %v", err)
	}

	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return b64, nil
}

func populateTitles(photos []*Photo, path string) ([]*Photo, error) {
	titles, err := parseTitles(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing titles: %v", err)
	}

	for _, photo := range photos {
		if title, ok := titles[photo.Path]; ok {
			log.Infof("Setting title %q for %s", title, photo.Path)
			photo.Title = title
		}
	}

	return photos, nil
}

func parseTitles(path string) (map[string]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	if strings.HasSuffix(path, ".json") {
		return parseTitlesInstagram(content)
	}

	return parseTitlesYaml(content)
}

func parseTitlesInstagram(content []byte) (map[string]string, error) {
	type image struct {
		URI string `json:"uri"`
	}
	type post struct {
		Media []image `json:"media"`
		Title string  `json:"title"`
	}
	type data []post

	var d data
	if err := json.Unmarshal(content, &d); err != nil {
		return nil, fmt.Errorf("error unmarshaling file: %v", err)
	}

	numberPrefix := regexp.MustCompile(`^\d+\.\s+`)

	titles := make(map[string]string)
	for _, post := range d {
		if post.Title == "" {
			continue
		}

		if len(post.Media) == 1 {
			titles[post.Media[0].URI] = post.Title
			continue
		}

		lines := strings.Split(post.Title, "\n")
		for i, media := range post.Media {
			if i >= len(lines) {
				log.Warnf("Not enough titles for %s", media.URI)
				break
			}
			// trim "media/posts/" from media.URI
			media.URI = strings.TrimPrefix(media.URI, "media/posts/")

			// trim number prefix if present
			lines[i] = numberPrefix.ReplaceAllString(lines[i], "")

			titles[media.URI] = lines[i]
		}
	}

	return titles, nil
}

func parseTitlesYaml(content []byte) (map[string]string, error) {
	var titles map[string]string
	if err := yaml.Unmarshal(content, &titles); err != nil {
		return nil, fmt.Errorf("error unmarshaling file: %v", err)
	}

	return titles, nil
}

func crc32sum(content []byte) string {
	hash := crc32.NewIEEE()
	if _, err := io.Copy(hash, bytes.NewReader(content)); err != nil {
		log.Errorf("error calculating CRC32 checksum: %v", err)
		return ""
	}

	return fmt.Sprintf("%x", hash.Sum32())
}
