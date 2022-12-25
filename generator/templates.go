package main

import (
	"log"
	"regexp"
	"sort"
	"strings"
	"text/template"

	i "github.com/nicksnyder/go-i18n/v2/i18n"
)

var fm = template.FuncMap{
	"config":     config,
	"enabled":    enabled,
	"i18n":       i18n,
	"year":       year,
	"alternates": alternates,
	"link":       link,
	"join":       join,
	"nextPage":   nextPage,
	"prevPage":   prevPage,
	"stripTags":  stripTags,
}

func config(key string) string {
	return cfg.GetString(key)
}

func enabled(key string) bool {
	return cfg.GetBool(key)
}

func i18n(key, lang string) string {
	localizer := i.NewLocalizer(bundle, lang)
	str, err := localizer.Localize(&i.LocalizeConfig{MessageID: key})
	if err != nil {
		log.Printf("error localizing %q: %v", key, err)
		return key
	}

	return str
}

func year(date string) string {
	if len(date) < 4 {
		return ""
	}

	return date[:4]
}

func alternates(data Data) []*MarkdownFile {
	if len(data.Alternates) > 0 {
		return data.Alternates
	}

	var result []*MarkdownFile

	for _, file := range data.All {
		if file.ID == data.File.ID {
			result = append(result, file)
		}
	}

	sort.Sort(ByLanguage(result))

	return result
}

var langSuffix = regexp.MustCompile(`_([a-z]{2}).(html|md)$`)

func link(path string, langs ...string) string {
	lang := ""
	if len(langs) > 0 {
		lang = langs[0]
	}

	link := path

	if langSuffix.MatchString(path) {
		match := langSuffix.FindStringSubmatch(path)
		lang = match[1]
		link = path[:len(path)-len(match[0])] + ".html"
	}

	if cfg.RemoveHTMLExtension {
		link = strings.TrimSuffix(
			strings.TrimSuffix(
				strings.TrimSuffix(link, "index.html"),
				".html",
			),
			"/",
		)
	}

	if lang != "" && lang != cfg.DefaultLanguage {
		link += "?lang=" + lang
	}

	return cfg.RootURL + link
}

func join(sep string, items []string) string {
	return strings.Join(items, sep)
}

func prevPage(data Data) (prev *MarkdownFile) {
	// technically, it get's the NEXT page
	// from the list of all pages SORTED by created date (descending)
	// but chronologically, it's the PREVIOUS page
	prev = nil

	var (
		i    int
		file *MarkdownFile
	)

	for i, file = range data.All {
		if file.Path == data.File.Path { // searching for the current page
			break
		}
	}

	for _, file := range data.All[i+1:] {
		if file.ID == data.File.ID { // skipping same pages in different languages
			continue
		}
		if file.Language == data.File.Language { // first page in the same language
			prev = file
			break
		}
	}

	return
}

func nextPage(data Data) (next *MarkdownFile) {
	// technically, it get's the PREVIOUS page
	// from the list of all pages SORTED by created date (descending)
	// but chronologically, it's the NEXT page
	next = nil

	for _, file := range data.All {
		if file.Path == data.File.Path { // searching for the current page
			break // this is the most recent page, so there's no next page
		}
		if file.Language == data.File.Language {
			// last seen page in the same language
			next = file
		}
	}

	return
}

var htmlTagRegexp = regexp.MustCompile("<[^>]*>")

func stripTags(html string) string {
	return htmlTagRegexp.ReplaceAllString(string(html), "")
}
