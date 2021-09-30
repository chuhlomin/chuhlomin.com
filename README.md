# micro

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

Use [fswatch](https://github.com/emcrisostomo/fswatch) to update static site
on every file change.

```bash
# brew install fswatch
fswatch -or -e "output" -e ".git" . | xargs -n1 genblog

# ⬆️ same as:
fswatch --one-per-batch --recursive --exclude="output" --exclude=".git" . | xargs -n1 genblog
```
