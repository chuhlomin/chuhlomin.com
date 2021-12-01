# micro

[![main](https://github.com/chuhlomin/micro/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/chuhlomin/micro/actions/workflows/main.yml)

Microblog, stored as Markdown files in Git repository.

GitHub Actions uses [genblog](https://github.com/chuhlomin/genblog) Go app that generates static site using `templates` directory.

Powers https://chuhlomin.com/blog/index.html

## Local setup

Set environment variables values from `.env` file.
You may use [alias e](https://github.com/chuhlomin/e) for that.

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

### Option 1. Open file in the browser

```
open output/index.html
```

⚠️ `?lang=ru` links will not work in this case.

### Option 1. [Docker](https://www.docker.com)

Serve output folder locally via Nginx using docker-compose:

```bash
docker-compose up -d nginx
# or
make run-docker
```

Open http://127.0.0.1:8080

### Option 2. [Podman](https://podman.io)

Or use [podman](https://podman.io) to build and run the image:

```bash
podman build -t micro:local .
podman run --name micro -p 8080:80 micro:local
# or
make build-podman
make run-podman
```

Open http://127.0.0.1:8080

### Option 3. [Caddy](https://caddyserver.com)

Add to `/etc/hosts` file:

```
127.0.0.1       local.chuhlomin.com
```

Install Caddy:

```bash
brew install caddy
```

Run Caddy from repository root:

```bash
caddy run
```

Open https://local.chuhlomin.com/blog/index.html
