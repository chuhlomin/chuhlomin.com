# homepage

[![main](https://github.com/chuhlomin/homepage/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/chuhlomin/homepage/actions/workflows/main.yml)

Powers https://chuhlomin.com

## Directory structure

```bash
.
├── content
│   ├── blog
│   │   ├── 2020
│   │   │   └── castty.md
│   │   ├── 2021
│   │   ├── 2022
│   │   ├── 2023
│   │   └── index.gohtml
│   ├── experiments
│   ├── fonts
│   ├── img
│   ├── styles.css
│   └── index.gohtml
├── generator # Go app
├── i18n      # translations
└── templates # Go templates used by generator
```

## Local development

Pre-requisites:

- [Go](https://go.dev/doc/install)
- [Caddy](https://caddyserver.com/docs/install)

```bash
make build
# will create `output` directory with generated static files

# in a separate terminal
carry run
```

Open https://local.chuhlomin.com

See [main.go](generator/main.go) for more configuration options.
