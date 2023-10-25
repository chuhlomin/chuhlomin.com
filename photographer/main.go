// A small go app to handle photos
// It will
// - generate thumbnails for all photos
// - upload photos and thumbnails to R2 Cloudflare storage if they don't exist yet
// - generate BlurHash for each photo
// - keep track of them in a photos.yaml file that will be used by the Generator to generate the website
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"
)

// Photo struct for items in photos.yml file
// Must be in sync with generator/photos.go
type Photo struct {
	Path     string
	Blurhash string `yaml:"blurhash,omitempty"`
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

	photos, err := loadPhotos(cfg.YamlFile)
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

	toAdd, toDelete := diff(photos, files)

	for _, file := range toAdd {
		photos = append(photos, Photo{
			Path: file,
		})

		content, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		log.Infof("Uploading %s", file)
		if err = r2.Upload(ctx, file, content); err != nil {
			return fmt.Errorf("error uploading file: %v", err)
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

	// save photos.yml file
	if err = savePhotos(cfg.YamlFile, photos); err != nil {
		return fmt.Errorf("error saving photos: %v", err)
	}

	return nil
}

func loadPhotos(path string) ([]Photo, error) {
	// read photos.yml file
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var photos []Photo
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
		if filepath.Ext(path) != ".jpg" {
			return nil
		}

		// trim the dir prefix
		path = strings.TrimPrefix(path, dir)
		files = append(files, path)
		return nil
	})

	return files, err
}

func diff(photos []Photo, files []string) ([]string, []string) {
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

func containsPhoto(arr []Photo, needle string) bool {
	for _, item := range arr {
		if item.Path == needle {
			return true
		}
	}

	return false
}

func savePhotos(path string, photos []Photo) error {
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
