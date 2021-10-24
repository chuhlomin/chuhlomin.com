# micro

[![main](https://github.com/chuhlomin/micro/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/chuhlomin/micro/actions/workflows/main.yml)

## Local setup

Use [e](https://github.com/chuhlomin/e) to get values from `.env` file into
environment:

```bash
e
```

Use [genblog](https://github.com/chuhlomin/genblog) binary to generate static
site from Markdown files and templates. Also copy files from `static` directory
to `output`

```bash
genblog
cp -R static/ output/

# same as:
make build
```

Use [fswatch](https://github.com/emcrisostomo/fswatch) to update the site on every file change.

```bash
# brew install fswatch
fswatch -or -e "output" -e ".git" . | xargs -n1 sh -c "genblog; cp -R static/ output/"

# same as:
fswatch --one-per-batch --recursive --exclude="output" --exclude=".git" . | xargs -n1 sh -c "genblog; cp -R static/ output/"

# same as:
make watch
```

Serve output folder locally via Nginx using docker-compose:

```bash
docker-compose up -d nginx
# or
make run-docker
```

Open http://127.0.0.1:8080
