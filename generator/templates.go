package main

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/log"
	i "github.com/nicksnyder/go-i18n/v2/i18n"
)

var fm = template.FuncMap{
	"config":     config,
	"enabled":    enabled,
	"i18n":       i18n,
	"date":       date,
	"alternates": alternates,
	"link":       link,
	"join":       join,
	"mkslice":    mkslice,
	"append":     appendSlice,
	"sort":       sortSlice,
	"uniq":       uniq,
	"nextPage":   nextPage,
	"prevPage":   prevPage,
	"stripTags":  stripTags,
	"hasSuffix":  strings.HasSuffix,
	"ts":         func() string { return ts.Format(time.RFC3339) },
	"jsonify":    jsonify,
	"divide":     func(a, b int) int { return a / b },
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
		log.Errorf("error localizing %q: %v", key, err)
		return key
	}

	return str
}

func date(date, format string) string {
	var t time.Time
	if date == "" {
		// use current time in UTC
		t = time.Now().UTC()
	} else {
		// parse date
		var err error
		t, err = time.Parse("2006-01-02", date)
		if err != nil {
			log.Errorf("error parsing date %q: %v", date, err)
			return date
		}
	}

	if format == "" {
		format = "2006-01-02"
	}

	return t.Format(format)
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

	return link
}

func join(sep string, items []string) string {
	return strings.Join(items, sep)
}

func mkslice() []string {
	return []string{}
}

func appendSlice(item string, items []string) []string {
	return append(items, item)
}

func sortSlice(items []string) []string {
	sort.Strings(items)
	return items
}

func uniq(items []string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, item := range items {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
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

	for i, file = range data.AllSorted {
		if file.Path == data.File.Path { // searching for the current page
			break
		}
	}

	for _, file := range data.AllSorted[i+1:] {
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

	for _, file := range data.AllSorted {
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

func jsonify(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		log.Errorf("error marshaling data: %v", err)
		return ""
	}
	return string(b)
}
