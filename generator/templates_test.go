package main

import (
	"strconv"
	"testing"
)

func TestLink(t *testing.T) {
	tt := []struct {
		cfg  *Config
		in   string
		lang string
		out  string
	}{
		{
			in:  "index.html",
			out: "index.html",
		},
		{
			cfg: &Config{
				RootURL: "https://example.com/",
			},
			in:  "index.html",
			out: "https://example.com/index.html",
		},
		{
			in:  "index_ru.html",
			out: "index.html?lang=ru",
		},
		{
			in:  "blog/index_ru.html",
			out: "blog/index.html?lang=ru",
		},
		{
			cfg: &Config{
				RemoveHTMLExtension: true,
			},
			in:  "index.html",
			out: "",
		},
		{
			cfg: &Config{
				RemoveHTMLExtension: true,
			},
			in:   "index.html",
			lang: "ru",
			out:  "?lang=ru",
		},
		{
			cfg: &Config{
				RemoveHTMLExtension: true,
			},
			in:  "blog/index.html",
			out: "blog",
		},
		{
			cfg: &Config{
				RemoveHTMLExtension: true,
			},
			in:  "blog/index_ru.html",
			out: "blog?lang=ru",
		},
		{
			cfg: &Config{
				RemoveHTMLExtension: true,
			},
			in:  "blog/2022/calendar_ru.html",
			out: "blog/2022/calendar?lang=ru",
		},
		{
			cfg: &Config{
				RemoveHTMLExtension: true,
				DefaultLanguage:     "en",
			},
			in:   "blog/2022/calendar.html",
			lang: "en",
			out:  "blog/2022/calendar",
		},
	}

	for ti, tc := range tt {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			if tc.cfg != nil {
				cfg = *tc.cfg
			} else {
				cfg = Config{}
			}

			out := link(tc.in, tc.lang)
			if out != tc.out {
				t.Errorf("expected %s, got %s", tc.out, out)
			}
		})
	}
}
