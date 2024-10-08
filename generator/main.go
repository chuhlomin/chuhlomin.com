// Generator is a static site generator for my site.
package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/log"
	flags "github.com/jessevdk/go-flags"
	"github.com/meilisearch/meilisearch-go"
	i "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Config is a struct that holds all configuration options
type Config struct {
	ContentDirectory    string        `env:"CONTENT_DIR" long:"content" description:"content directory" default:"content"`
	TemplatesDirectory  string        `env:"TEMPLATES_DIR" long:"templates" description:"templates directory" default:"templates"`
	OutputDirectory     string        `env:"OUTPUT_DIR" long:"output" description:"output directory" default:"output"`
	TempDirectory       string        `env:"TEMP_DIR" long:"temp" description:"temp directory" default:""`
	DefaultLanguage     string        `env:"DEFAULT_LANG" long:"default-lang" description:"default language" default:"en"`
	RootURL             string        `env:"ROOT_URL" long:"root-url" description:"root url" default:"https://local.chuhlomin.com"`
	CommentsSiteID      string        `env:"COMMENTS_SITE_ID" long:"comments-site-id" description:"comments site id"`
	CacheDirectory      string        `env:"CACHE_DIR" long:"cache-dir" description:"cache directory" default:"cache"`
	I18NDirectory       string        `env:"I18N_DIR" long:"i18n-dir" description:"i18n directory" default:"i18n"`
	SearchHost          string        `env:"SEARCH_HOST" long:"search-host" description:"search host" default:"https://local.chuhlomin.com/blog/search"`
	SearchMasterKey     string        `env:"SEARCH_MASTER_KEY" long:"search-master-key" description:"search master key, used to create index"`
	SearchAPIKey        string        `env:"SEARCH_API_KEY" long:"search-api-key" description:"search api key, used on frontend to search"`
	ThumbPath           string        `env:"THUMB_PATH" long:"thumb-path" description:"path to thumbnails" default:"/img/thumbs/"`
	OpenGraphCacheFile  string        `env:"OPENGRAPH_CACHE_FILE" long:"opengraph-cache-file" description:"opengraph cache file" default:"cache.yml"`
	PhotosDomain        string        `env:"PHOTOS_DOMAIN" long:"photos-domain" description:"photos domain" default:"https://photos.chuhlomin.com/"`
	FilesChannelSize    int           `env:"FILES_CHANNEL_SIZE" long:"files-channel-size" description:"size of file channel" default:"100"`
	ImagesChannelSize   int           `env:"IMAGES_CHANNEL_SIZE" long:"images-channel-size" description:"size of images channel" default:"100"`
	SearchTimeout       time.Duration `env:"SEARCH_TIMEOUT" long:"search-timeout" description:"search timeout" default:"5s"`
	NumWorkers          int           `env:"NUM_WORKERS" long:"workers" description:"number of workers" default:"4"`
	ThumbMaxWidth       int           `env:"THUMB_MAX_WIDTH" long:"thumb-max-width" description:"max width of thumbnails" default:"140"`
	ThumbMaxHeight      int           `env:"THUMB_MAX_HEIGHT" long:"thumb-max-height" description:"max height of thumbnails" default:"140"`
	OpenGraphTimeout    time.Duration `env:"OPENGRAPH_TIMEOUT" long:"opengraph-timeout" description:"opengraph timeout" default:"5s"`
	Debug               bool          `env:"DEBUG" long:"debug" description:"debug mode"`
	RemoveHTMLExtension bool          `env:"REMOVE_HTML_EXT" long:"remove-html-ext" description:"remove .html extension from urls"`
	CommentsEnabled     bool          `env:"COMMENTS_ENABLED" long:"comments-enabled" description:"enable comments"`
	ShowDrafts          bool          `env:"SHOW_DRAFTS" long:"show-drafts" description:"show drafts"`
	SearchEnabled       bool          `env:"SEARCH_ENABLED" long:"search-enabled" description:"enable search"`
}

// GetString returns the value of the environment variable named by the key.
// If the variable is not present, GetString returns empty string.
// Used in `config` template function to access config values.
func (c Config) GetString(key string) string {
	// use reflect to get the value of the key
	v := reflect.ValueOf(c)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() != reflect.String {
			continue
		}

		if v.Type().Field(i).Name == key {
			return v.Field(i).String()
		}
	}
	return ""
}

// GetBool returns the value of the environment variable named by the key.
// If the variable is not present, GetBool returns `false`.
// Used in `enabled` template function to access config values.
func (c Config) GetBool(key string) bool {
	// use reflect to get the value of the key
	v := reflect.ValueOf(c)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() != reflect.Bool {
			continue
		}

		if v.Type().Field(i).Name == key {
			return v.Field(i).Bool()
		}
	}
	return false
}

var (
	ts     time.Time // timestamp used to measure execution time
	cfg    Config    // global config
	bundle *i.Bundle // used in templates/i18n to get translated strings
)

func main() {
	log.Info("Starting")
	ts = time.Now()

	if err := run(ts); err != nil {
		log.Fatal(err)
	}

	log.Infof("Finished in %v", time.Now().Sub(ts))
}

func run(ts time.Time) error {
	_, err := flags.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("Error parsing flags: %v", err)
	}

	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if err := initBundle(); err != nil {
		return fmt.Errorf("Error initializing i18n bundle: %v", err)
	}

	og, err := newOpenGraphClient(cfg.OpenGraphTimeout, cfg.OpenGraphCacheFile)
	if err != nil {
		return fmt.Errorf("Error creating opengraph client: %v", err)
	}
	defer og.Save(cfg.OpenGraphCacheFile)

	var searchClient *meilisearch.Client
	if cfg.SearchEnabled {
		searchClient = meilisearch.NewClient(meilisearch.ClientConfig{
			Host:    cfg.SearchHost,
			APIKey:  cfg.SearchMasterKey,
			Timeout: cfg.SearchTimeout,
		})
	}

	generator, err := NewGenerator(og, searchClient)
	if err != nil {
		return fmt.Errorf("Error creating generator: %v", err)
	}

	return generator.Run(ts)
}

func initBundle() error {
	if cfg.DefaultLanguage == "" {
		cfg.DefaultLanguage = "en"
	}

	lang, err := language.Parse(cfg.DefaultLanguage)
	if err != nil {
		return fmt.Errorf("parse language %q", cfg.DefaultLanguage)
	}
	bundle = i.NewBundle(lang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	return nil
}
