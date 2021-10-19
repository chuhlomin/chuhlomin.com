# micro

[![main](https://github.com/chuhlomin/micro/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/chuhlomin/micro/actions/workflows/main.yml)

## Local setup

Use [e](https://github.com/chuhlomin/e) to get values from `.env` file into
environment:

```bash
e
```

Use [genblog](https://github.com/chuhlomin/genblog) binary to generate static
site from Markdown files and templates.

```bash
genblog
```

Use [fswatch](https://github.com/emcrisostomo/fswatch) to update the site on every file change.

```bash
# brew install fswatch
fswatch -or -e "output" -e ".git" . | xargs -n1 sh -c "genblog; cp -R static/ output/"

# ⬆️ same as:
fswatch --one-per-batch --recursive --exclude="output" --exclude=".git" . | xargs -n1 sh -c "genblog; cp -R static/ output/"
```

Serve output folder locally via nginx using docker-compose:

```bash
docker-compose up -d nginx
```

Open http://127.0.0.1:8080
